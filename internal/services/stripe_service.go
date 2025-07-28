package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"peruccii/site-vigia-be/db"
	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/repository"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeService interface {
	CreateCustomer(ctx context.Context, name, email string) (*stripe.Customer, error)
	CreateSubscriptionCheckout(ctx context.Context, req *CreateSubscriptionRequest) (*stripe.CheckoutSession, error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
	ProcessWebhook(ctx context.Context, payload []byte, signature string) error
	GetSubscription(subscriptionID string) (*stripe.Subscription, error)
}

type StripeProvider struct {
	secretKey        string
	webhookSecret    string
	paymentRepo      repository.PaymentRepository
	subscriptionRepo repository.SubscriptionRepository
	planRepo         repository.PlanRepository
}

func NewStripeService(
	paymentRepo repository.PaymentRepository,
	subscriptionRepo repository.SubscriptionRepository,
	planRepo repository.PlanRepository,
) StripeService {
	return &StripeProvider{
		secretKey:        os.Getenv("STRIPE_SECRET_KEY"),
		webhookSecret:    os.Getenv("STRIPE_WEBHOOK_SECRET"),
		paymentRepo:      paymentRepo,
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
	}
}

type CreateSubscriptionRequest struct {
	UserID     uuid.UUID `json:"user_id"`
	PlanID     uuid.UUID `json:"plan_id"`
	UserName   string    `json:"user_name"`
	UserEmail  string    `json:"user_email"`
	SuccessURL string    `json:"success_url"`
	CancelURL  string    `json:"cancel_url"`
}

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSucceeded PaymentStatus = "succeeded"
	PaymentFailed    PaymentStatus = "failed"
	PaymentCanceled  PaymentStatus = "canceled"
	PaymentRefunded  PaymentStatus = "refunded"
)

// done
func (s *StripeProvider) CreateCustomer(ctx context.Context, name, email string) (*stripe.Customer, error) {
	stripe.Key = s.secretKey

	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
		Metadata: map[string]string{
			"source": "site-vigia",
		},
	}

	result, err := customer.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return result, nil
}

// TODO: Plan PriceMonthly must be a price ID stripe
func (s *StripeProvider) CreateSubscriptionCheckout(ctx context.Context, req *CreateSubscriptionRequest) (*stripe.CheckoutSession, error) {
	stripe.Key = s.secretKey

	plan, err := s.planRepo.GetPlanByID(ctx, req.PlanID)
	if err != nil {
		return nil, fmt.Errorf("plano não encontrado: %w", err)
	}

	if *&plan.PriceMonthly == "" {
		return nil, fmt.Errorf("plano %s não tem integração com Stripe configurada", plan.Name)
	}

	successURL := req.SuccessURL
	if successURL == "" {
		successURL = os.Getenv("STRIPE_SUCCESS_URL")
	}
	successURL += "?session_id={CHECKOUT_SESSION_ID}&plan=" + plan.Name

	cancelURL := req.CancelURL
	if cancelURL == "" {
		cancelURL = os.Getenv("STRIPE_CANCEL_URL")
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL:    stripe.String(successURL),
		CancelURL:     stripe.String(cancelURL),
		Mode:          stripe.String(stripe.CheckoutSessionModeSubscription),
		CustomerEmail: stripe.String(req.UserEmail),

		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"pix",
		}),

		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(plan.PriceMonthly),
				Quantity: stripe.Int64(1),
			},
		},

		Metadata: map[string]string{
			"user_id":   req.UserID.String(),
			"plan_id":   req.PlanID.String(),
			"plan_name": plan.Name,
			"source":    "site-vigia",
		},

		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": req.UserID.String(),
				"plan_id": req.PlanID.String(),
			},
		},

		BillingAddressCollection: stripe.String("required"),
	}

	result, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar sessão de checkout: %w", err)
	}
	log.Printf("Checkout session created: %s for user %s, plan %s",
		result.ID, req.UserID, plan.Name)

	return result, nil
}

func (s *StripeProvider) CancelSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	stripe.Key = s.secretKey

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := subscription.Update(subscriptionID.String(), params)
	if err != nil {
		return fmt.Errorf("erro ao cancelar assinatura: %w", err)
	}

	input := db.UpdateSubscriptionStatusParams{
		ID:     subscriptionID,
		Status: "canceled",
	}

	err = s.subscriptionRepo.UpdateSubscriptionStatus(ctx, input)
	if err != nil {
		log.Printf("Erro ao atualizar status de cancelamento no banco: %v", err)
	}

	return nil
}

func (s *StripeProvider) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	stripe.Key = s.secretKey
	return subscription.Get(subscriptionID, nil)
}

func (s *StripeProvider) ProcessWebhook(ctx context.Context, payload []byte, signature string) error {
	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	switch event.Type {
	case "checkout.session.completed":
		return s.handleCheckoutCompleted(ctx, event)
	case "customer.subscription.created":
		return s.handleSubscriptionCreated(ctx, event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	case "invoice.payment_succeeded":
		return s.handlePaymentSucceeded(ctx, event)
	case "invoice.payment_failed":
		return s.handlePaymentFailed(ctx, event)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return nil
}

func (s *StripeProvider) handleCheckoutCompleted(ctx context.Context, event stripe.Event) error {
	var session stripe.CheckoutSession

	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return fmt.Errorf("erro ao fazer parse do evento: %w", err)
	}

	userID, err := uuid.Parse(session.Metadata["user_id"])
	if err != nil {
		return fmt.Errorf("invalid user_id in metadata: %w", err)
	}

	planID, err := uuid.Parse(session.Metadata["plan_id"])
	if err != nil {
		return fmt.Errorf("invalid plan_id in metadata: %w", err)
	}

	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	paymentReq := &dto.CreatePaymentRequest{
		UserID:            userID,
		PlanID:            planID,
		StripeSessionID:   &session.ID,
		Amount:            plan.PriceMonthly,
		Status:            string(PaymentPending),
		PaymentMethodType: "checkout_session",
	}

	err = s.paymentRepo.Create(ctx, paymentReq)
	if err != nil {
		return fmt.Errorf("error creating payment record: %w", err)
	}

	log.Printf("Checkout completed for user %s, plan %s, session %s",
		userID, plan.Name, session.ID)

	return nil
}

func (s *StripeProvider) handlePaymentSucceeded(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("error parsing invoice  %w", err)
	}

	payment, err := s.paymentRepo.FindByStripeSubscriptionID(ctx, invoice.Customer.Subscriptions.Data[0].ID)
	if err != nil {
		return fmt.Errorf("payment not found for subscription %s: %w", invoice.Subscription.ID, err)
	}

	now := time.Now()
	err = s.paymentRepo.UpdateStatus(ctx, payment.ID, string(PaymentSucceeded), nil, &now)
	if err != nil {
		return fmt.Errorf("error updating payment status: %w", err)
	}

	log.Printf("Payment succeeded: %s for user %s", payment.ID, payment.UserID)
	return nil
}

func (s *StripeProvider) handlePaymentFailed(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("error parsing invoice: %w", err)
	}
	payment, err := s.paymentRepo.FindByStripeSubscriptionID(ctx, invoice.Subscription.ID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	failureReason := "Invoice payment failed"
	if invoice.LastPaymentError != nil {
		failureReason = invoice.LastPaymentError.Message
	}

	err = s.paymentRepo.UpdateStatus(ctx, payment.ID, string(PaymentFailed), &failureReason, nil)
	if err != nil {
		return fmt.Errorf("error updating payment status: %w", err)
	}

	log.Printf("Payment failed: %s for user %s, reason: %s",
		payment.ID, payment.UserID, failureReason)

	return nil
}

func (s *StripeProvider) handleSubscriptionCreated(ctx context.Context, event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return fmt.Errorf("error parsing subscription: %w", err)
	}
	userID, err := uuid.Parse(subscription.Metadata["user_id"])
	if err != nil {
		return fmt.Errorf("invalid user_id in metadata: %w", err)
	}

	planID, err := uuid.Parse(subscription.Metadata["plan_id"])
	if err != nil {
		return fmt.Errorf("invalid plan_id in metadata: %w", err)
	}

	subscriptionReq := &dto.CreateSubscriptionRequest{
		UserID:               userID,
		PlanID:               planID,
		StripeCustomerID:     subscription.Customer.ID,
		StripeSubscriptionID: subscription.ID,
		Status:               subscription.Status,
		CurrentPeriodStart:   time.Unix(subscription.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(subscription.CurrentPeriodEnd, 0),
		CancelAtPeriodEnd:    subscription.CancelAtPeriodEnd,
	}

	_, err = s.subscriptionRepo.CreateOrUpdate(ctx, subscriptionReq)
	if err != nil {
		return fmt.Errorf("error creating subscription: %w", err)
	}

	log.Printf("Subscription created: %s for user %s", subscription.ID, userID)
	return nil
}

func (s *StripeProvider) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) error {
	// Similar ao created, mas para updates
	return s.handleSubscriptionCreated(ctx, event)
}

func (s *StripeProvider) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return fmt.Errorf("error parsing subscription: %w", err)
	}

	input := db.UpdateSubscriptionStatusParams{
		ID:     subscription.ID,
		Status: "deleted",
	}

	err := s.subscriptionRepo.UpdateSubscriptionStatus(ctx, input)
	if err != nil {
		return fmt.Errorf("error updating subscription status: %w", err)
	}

	log.Printf("Subscription deleted: %s", subscription.ID)
	return nil
}
