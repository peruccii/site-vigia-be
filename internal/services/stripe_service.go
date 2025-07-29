package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"peruccii/site-vigia-be/db"
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
	CancelSubscription(ctx context.Context, subscriptionID uuid.UUID) error
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
	PlanID     int32     `json:"plan_id"`
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

	if plan.StripePriceID == nil || *plan.StripePriceID == "" {
		return nil, fmt.Errorf("plano %s não tem um Price ID do Stripe configurado", plan.Name)
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
				Price:    stripe.String(*plan.StripePriceID),
				Quantity: stripe.Int64(1),
			},
		},

		Metadata: map[string]string{
			"user_id":   req.UserID.String(),
			"plan_id":   strconv.Itoa(int(req.PlanID)),
			"plan_name": plan.Name,
			"source":    "site-vigia",
		},

		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": req.UserID.String(),
				"plan_id": strconv.Itoa(int(req.PlanID)),
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
	case "invoice.payment_succeeded":
		return s.handlePaymentSucceeded(ctx, event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return nil
}

func (s *StripeProvider) handleCheckoutCompleted(ctx context.Context, event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return fmt.Errorf("erro ao fazer parse da sessão do evento: %w", err)
	}

	if session.Mode != stripe.CheckoutSessionModeSubscription {
		return nil
	}

	userID, err := uuid.Parse(session.Metadata["user_id"])
	if err != nil {
		return fmt.Errorf("invalid user_id in metadata: %w", err)
	}

	// converting string to int64
	planID, err := strconv.ParseInt(session.Metadata["plan_id"], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid plan_id in metadata: %w", err)
	}

	stripeSubID := session.Subscription.ID
	if stripeSubID == "" {
		return fmt.Errorf("stripe_subscription_id não encontrado na sessão de checkout")
	}

	stripeSub, err := s.GetSubscription(stripeSubID)
	if err != nil {
		return fmt.Errorf("erro ao buscar assinatura no Stripe %s: %w", stripeSubID, err)
	}

	subParams := db.CreateSubscriptionParams{
		UserID:               userID,
		PlanID:               int32(planID),
		Status:               string(stripeSub.Status),
		StripeSubscriptionID: &stripeSubID,
		CurrentPeriodEndsAt:  time.Unix(stripeSub.EndedAt, 0),
	}

	err = s.subscriptionRepo.CreateSubscription(ctx, subParams)
	if err != nil {
		// Adicionar lógica para lidar com o caso de a assinatura já existir (idempotência)
		return fmt.Errorf("erro ao criar registro de assinatura: %w", err)
	}

	log.Printf("Assinatura %s criada com sucesso para o usuário %s", stripeSubID, userID)
	return nil
}

func (s *StripeProvider) handlePaymentSucceeded(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("erro ao fazer parse da fatura (invoice): %w", err)
	}

	stripeSubID := invoice.Subscription.ID
	if stripeSubID == "" {
		return nil
	}

	subscription, err := s.subscriptionRepo.GetSubscriptionByStripeSbId(ctx, stripeSubID)
	if err != nil {
		return fmt.Errorf("assinatura com stripe_id %s não encontrada no banco: %w", stripeSubID, err)
	}

	paymentParams := db.CreatePaymentParams{
		UserID:                subscription.UserID,
		SubscriptionID:        subscription.ID,
		StripePaymentIntentID: invoice.PaymentIntent.ID,
		StripeInvoiceID:       &invoice.ID,
		AmountCents:           int(invoice.AmountPaid), // Stripe retorna em centavos
		Currency:              string(invoice.Currency),
		Status:                string(PaymentSucceeded),
		PaymentMethod:         string(invoice.PaymentSettings.PaymentMethodTypes[0]), // Simplificação
	}
	// paidAt := time.Unix(invoice.StatusTransitions.PaidAt, 0)
	// paymentParams.PaidAt = sql.NullTime{Time: paidAt, Valid: true}

	err = s.paymentRepo.Create(ctx, paymentParams)
	if err != nil {
		return fmt.Errorf("erro ao criar registro de pagamento: %w", err)
	}

	stripeSub, err := s.GetSubscription(stripeSubID)
	if err != nil {
		return fmt.Errorf("erro ao buscar assinatura no Stripe para atualização: %w", err)
	}

	updateSubParams := db.UpdateSubscriptionPeriodParams{
		ID:                  subscription.ID,
		Status:              string(stripeSub.Status),
		CurrentPeriodEndsAt: stripeSub.CurrentPeriodEnd,
	}
	err = s.subscriptionRepo.UpdatePeriod(ctx, updateSubParams)
	if err != nil {
		return fmt.Errorf("erro ao atualizar período da assinatura: %w", err)
	}

	log.Printf("Pagamento %s registrado e assinatura %s atualizada", invoice.PaymentIntent.ID, stripeSubID)
	return nil
}

func (s *StripeProvider) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		return fmt.Errorf("error parsing subscription: %w", err)
	}

	subscription, err := s.subscriptionRepo.GetSubscriptionByStripeSbId(ctx, stripeSub.ID)
	if err != nil {
		return fmt.Errorf("assinatura com stripe_id %s não encontrada para atualização: %w", stripeSub.ID, err)
	}

	newPriceID := stripeSub.Items.Data[0].Price.ID
	newPlan, err := s.planRepo.FindByStripePriceID(ctx, newPriceID)
	if err != nil {
		return fmt.Errorf("novo plano com price_id %s não encontrado: %w", newPriceID, err)
	}

	// Atualizar o registro no banco
	updateParams := db.UpdateSubscriptionPlanParams{
		ID:                  subscription.ID,
		PlanID:              newPlan.ID,
		Status:              string(stripeSub.Status),
		CurrentPeriodEndsAt: stripeSub.CurrentPeriodEnd,
	}

	err = s.subscriptionRepo.UpdatePlan(ctx, updateParams)
	if err != nil {
		return fmt.Errorf("erro ao atualizar plano da assinatura: %w", err)
	}

	log.Printf("Assinatura %s atualizada para o plano %s", stripeSub.ID, newPlan.Name)
	return nil
}

func (s *StripeProvider) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		return fmt.Errorf("error parsing subscription: %w", err)
	}

	err := s.subscriptionRepo.UpdateStatusByStripeID(ctx, stripeSub.ID, string(stripeSub.Status))
	if err != nil {
		return fmt.Errorf("erro ao atualizar status de cancelamento da assinatura %s: %w", stripeSub.ID, err)
	}

	log.Printf("Assinatura %s marcada como '%s' no banco de dados.", stripeSub.ID, stripeSub.Status)
	return nil
}
