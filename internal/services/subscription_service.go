package services

import (
	"context"
	"database/sql"
	"fmt"

	"peruccii/site-vigia-be/db"
	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/repository"

	"github.com/go-playground/validator/v10"
)

type subscriptionService struct {
	repo      *repository.SubscriptionRepository
	validator *validator.Validate
}

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, input dto.CreateSubscriptionRequest) error
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *subscriptionService {
	return &subscriptionService{repo: repo, validator: validator.New()}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, input dto.CreateSubscriptionRequest) error {
	if err := s.validator.Struct(input); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	existingSubscription, err := s.repo.GetSubscriptionByStripeSbId(ctx, input.Email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing subscription: %w", err)
	}

	if err == nil && existingSubscription.StripeSubscriptionID.Valid {
		return ErrUserAlreadyExists
	}

	subscription := db.CreateSubscriptionParams{
		UserID:               input.UserID,
		PlanID:               input.PlanID,
		Status:               input.Status,
		StripeSubscriptionID: input.StripeSubscriptionID,
		CurrentPeriodEndsAt:  input.CurrentPeriodEndsAt,
	}

	if err := s.repo.CreateSubscription(ctx, subscription); err != nil {
		return err
	}

	return nil
}
