package dto

import "github.com/google/uuid"

type CreatePaymentRequest struct {
	UserID            uuid.UUID `json:"user_id"`
	PlanID            uuid.UUID `json:"plan_id"`
	StripeSessionID   *string   `json:"stripe_session_id"`
	Amount            string    `json:"amount"`
	Currency          string    `json:"currency"`
	Status            string    `json:"status"`
	PaymentMethodType string    `json:"payment_method_type"`
}
