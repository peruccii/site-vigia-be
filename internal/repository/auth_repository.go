package repository

import (
	"context"

	"peruccii/site-vigia-be/db"
)

type AuthRepository struct {
	queries *db.Queries
}

func NewAuthRepository(queries *db.Queries) *AuthRepository {
	return &AuthRepository{
		queries: queries,
	}
}

func (r *AuthRepository) RegisterUser(ctx context.Context, input db.RegisterUserParams) error {
	return r.queries.RegisterUser(ctx, input)
}

func (r *AuthRepository) SignInUser(ctx context.Context, input db.SignInUserParams) (db.User, error) {
	return r.queries.SignInUser(ctx, input)
}
