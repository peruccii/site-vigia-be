package services

import (
	"context"
	"database/sql"
	"fmt"

	"peruccii/site-vigia-be/db"
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
	RecoverPassword(ctx context.Context, input dto.RecoverPasswordRequest) (string, error)
	ResetPassword(ctx context.Context, email string, input dto.ResetPasswordRequest) error
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
		return dto.SignInUserResponse{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	existingUser, err := s.userRepo.GetUserByEmail(ctx, input.Email)

	if err != nil && err != sql.ErrNoRows {
		return dto.SignInUserResponse{}, fmt.Errorf("failed to check existing user: %w", err)
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

func (s *authService) RecoverPassword(ctx context.Context, input dto.RecoverPasswordRequest) (string, error) {
	if err := s.validator.Struct(input); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	_, err := s.userRepo.GetUserByEmail(ctx, input.Email)

	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to check existing user: %w", err)
	}

	token, _ := utils.GenerateEmailToken(input.Email)

	SendEmail(input.Email, token)

	return "email sent", nil
}

func (s *authService) IsTokenValid(ctx context.Context, token string) (string, bool) {
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return "token is invalid", false
	}
	email, ok := claims["email"].(string)
	if !ok {
		return "email claim is missing or invalid", false
	}
	return email, true
}

func (s *authService) ResetPassword(ctx context.Context, email string, input dto.ResetPasswordRequest) error {
	if err := s.validator.Struct(input); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	user, _ := s.userRepo.GetUserByEmail(ctx, email)

	newPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("%w %v", ErrHashPassword, err)
	}

	if err := s.repo.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:           user.ID,
		PasswordHash: newPassword,
	}); err != nil {
		return fmt.Errorf("failed to update password %w", err)
	}
	return nil
}
