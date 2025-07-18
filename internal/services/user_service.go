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

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo      *repository.UserRepository
	validator *validator.Validate
}

type UserService interface {
	RegisterUser(ctx context.Context, input dto.RegisterUserRequest) error
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
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	exitingUser, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if err == nil && exitingUser.Email != "" {
		return ErrUserAlreadyExists
	}

	hashPassword, err := hashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrHashPassword, err)
	}

	user := db.RegisterUserParams{
		Name:            strings.TrimSpace(input.Name),
		Email:           strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash:    hashPassword,
		EmailVerifiedAt: sql.NullTime{Valid: false},
	}

	if err := s.repo.RegisterUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}
