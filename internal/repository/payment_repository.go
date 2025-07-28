package repository

import (
	"context"

	"peruccii/site-vigia-be/db"
)

type PaymentRepository interface {
	Create(ctx context.Context, input db.CreatePaymentParams) error
}

type paymentRepository struct {
	queries *db.Queries
}

func NewPaymentRepository(queries *db.Queries) PaymentRepository {
	return &paymentRepository{queries: queries}
}

func (r *paymentRepository) Create(ctx context.Context, input db.CreatePaymentParams) error {
	return r.queries.CreatePayment(ctx, input)
}
