package repository

import (
	"context"

	"peruccii/site-vigia-be/db"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

func (r *UserRepository) RegisterUser(ctx context.Context, input db.RegisterUserParams) error {
	return r.queries.RegisterUser(ctx, input)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.queries.GetUser(ctx, email)
}
