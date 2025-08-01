package repository

import (
	"context"

	"peruccii/site-vigia-be/db"
)

type SubscriptionRepository struct {
	queries *db.Queries
}

func NewSubscriptionRepository(queries *db.Queries) *SubscriptionRepository {
	return &SubscriptionRepository{
		queries: queries,
	}
}

func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, input db.CreateSubscriptionParams) error {
	return r.queries.CreateSubscription(ctx, input)
}

func (r *SubscriptionRepository) GetSubscriptionByStripeSbId(ctx context.Context, subscription_stripe_id string) (db.Subscription, error) {
	return r.queries.GetSubscriptionByStripeSbId(ctx, subscription_stripe_id)
}

func (r *SubscriptionRepository) UpdateSubscriptionStatus(ctx context.Context, input db.UpdateSubscriptionStatusParams) error {
	return r.queries.UpdateSubscriptionStatus(ctx, input)
}

func (r *SubscriptionRepository) UpdateSubscriptionPlan(ctx context.Context, input db.UpdateSubscriptionPlanParams) error {
	return r.queries.UpdateSubscriptionPlan(ctx, input)
}

func (r *SubscriptionRepository) UpdateStatusByStripeID(ctx context.Context, input db.UpdateStatusByStripeIDParams) error {
	return r.queries.UpdateStatusByStripeID(ctx, input)
}

func (r *SubscriptionRepository) UpdatePeriod(ctx context.Context, input db.UpdateSubscriptionPeriodParams) error {
	return r.queries.UpdateSubscriptionPeriod(ctx, input)
}
