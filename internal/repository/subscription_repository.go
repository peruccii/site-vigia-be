package repository

import (
	"context"
	"database/sql"

	"peruccii/site-vigia-be/db"
)

type SubscriptionRepository struct {
	queries *db.Queries
}

func NewSubscriptionRepository(dbConn *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		queries: db.New(dbConn),
	}
}

func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, input db.CreateSubscriptionParams) error {
	return r.queries.CreateSubscription(ctx, input)
}
