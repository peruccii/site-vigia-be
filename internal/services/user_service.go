package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"peruccii/site-vigia-be/db"
	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/repository"
	"peruccii/site-vigia-be/internal/utils"

	"github.com/go-playground/validator/v10"
)

type userService struct {
	repo      *repository.UserRepository
	validator *validator.Validate
}

type UserService interface {
	RegisterUser(ctx context.Context, input dto.RegisterUserRequest) error
	FindByEmail(ctx context.Context, email string) (db.User, error)
}

func NewUserService(repo *repository.UserRepository) *userService {
	return &userService{repo: repo, validator: validator.New()}
}

const (
	bcryptCost = 12
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrHashPassword      = errors.New("failed to hash password")
)

func (s *userService) RegisterUser(ctx context.Context, input dto.RegisterUserRequest) error {
	if err := s.validator.Struct(input); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err.Error())
	}

	exitingUser, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if err == nil && exitingUser.Email != "" {
		return ErrUserAlreadyExists
	}

	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrHashPassword, err)
	}

	user := db.RegisterUserParams{
		Name:            strings.TrimSpace(input.Name),
		Email:           strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash:    hashPassword,
		EmailVerifiedAt: nil,
	}

	if err := s.repo.RegisterUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) FindByEmail(ctx context.Context, email string) (db.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
