package services

import (
	"context"
	"database/sql"
	"fmt"

	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/repository"
	"peruccii/site-vigia-be/internal/utils"

	"github.com/go-playground/validator/v10"
)

type authService struct {
	repo      *repository.AuthRepository
	userRepo  *repository.UserRepository
	validator *validator.Validate
}

type AuthService interface {
	SignInUser(ctx context.Context, input dto.SignInUserRequest) (dto.SignInUserResponse, error)
}

func NewAuthService(repo *repository.AuthRepository, userRepo *repository.UserRepository) *authService {
	return &authService{
		repo:      repo,
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

func (s *authService) SignInUser(ctx context.Context, input dto.SignInUserRequest) (dto.SignInUserResponse, error) {
	if err := s.validator.Struct(input); err != nil {
		return dto.SignInUserResponse{}, fmt.Errorf("%wr: %var", ErrInvalidInput)
	}

	existingUser, err := s.userRepo.GetUserByEmail(ctx, input.Email)

	if err != nil && err != sql.ErrNoRows {
		return dto.SignInUserResponse{}, fmt.Errorf("failed to check existing user: %w", err)
	}

	if err == nil && existingUser.Email != "" {
		return dto.SignInUserResponse{}, ErrUserAlreadyExists
	}

	utils.VerifyPassword(input.Password, existingUser.PasswordHash)

	access_token, _ := utils.GenerateAccessToken(
		existingUser.ID.String(),
		existingUser.Name,
		existingUser.Email,
	)

	refresh_token, _ := utils.GenerateRefreshToken(existingUser.ID.String())

	return dto.SignInUserResponse{
		UserID:       existingUser.ID,
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		ExpiresIn:    10,
	}, nil
}
