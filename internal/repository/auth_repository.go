package repository

import (
	"context"

	"peruccii/site-vigia-be/db"
	"peruccii/site-vigia-be/internal/dto"
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

func (r *AuthRepository) SignInUser(ctx context.Context, input db.SignInUserParams) (dto.SignInUserResponse, error) {
	return r.queries.SignInUser(ctx, input)
}
