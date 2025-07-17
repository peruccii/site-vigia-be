package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	UserID               uuid.UUID  `json:"user_id" db:"user_id"`
	PlanID               int        `json:"plan_id" db:"plan_id"`
	Status               string     `json:"status" db:"status"` // active, cancelled, past_due
	StripeSubscriptionID *string    `json:"stripe_subscription_id" db:"stripe_subscription_id"`
	CurrentPeriodEndsAt  *time.Time `json:"current_period_ends_at" db:"current_period_ends_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}
