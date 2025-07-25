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

func (r *SubscriptionRepository) GetSubscriptionByStripeSbId(ctx context.Context, subscription_stripe_id *string) (db.Subscription, error) {
	return r.queries.GetSubscriptionByStripeSbId(ctx, subscription_stripe_id)
}
