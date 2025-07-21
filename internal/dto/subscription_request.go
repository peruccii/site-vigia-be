package dto

type CreateSubscriptionRequest struct {
	UserID               string `json:"user_id" validate:"required,uuid"`
	PlanID               int    `json:"plan_id" validate:"required,min=1,max=100"`
	Status               string `json:"status" validate:"required,oneof=active suspended"`
	StripeSubscriptionID string `json:"stripe_subscription_id" validate:"required,uuid"`
	CurrentPeriodEndsAt  string `json:"current_period_ends_at" validate:"required,datetime"`
}
