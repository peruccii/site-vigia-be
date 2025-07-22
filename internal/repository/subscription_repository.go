package repository

import (
	"context"
	"database/sql"

	"peruccii/site-vigia-be/db"

	"github.com/google/uuid"
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

func (r *SubscriptionRepository) GetSubscriptionByStripeSbId(ctx context.Context, subscription_stripe_id uuid.UUID) (db.Subscription, error) {
	return r.queries.GetSubscriptionByStripeSbId(ctx, subscription_stripe_id)
}
