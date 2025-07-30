package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	UserID               uuid.UUID  `json:"user_id" validate:"required,uuid"`
	PlanID               int32      `json:"plan_id" validate:"required,min=1,max=100"`
	Status               string     `json:"status" validate:"required,oneof=active suspended"`
	StripeSubscriptionID string    `json:"stripe_subscription_id" validate:"required,uuid"`
	CurrentPeriodEndsAt  time.Time `json:"current_period_ends_at" validate:"required,datetime"`
}
