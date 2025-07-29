package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"peruccii/site-vigia-be/db"
	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/repository"

	"github.com/go-playground/validator/v10"
)

type planService struct {
	repo      *repository.PlanRepository
	validator *validator.Validate
}

type PlanService interface {
	CreatePlan(ctx context.Context, input dto.CreatePlanRequest) error
}

func NewPlanService(repo *repository.PlanRepository) *planService {
	return &planService{repo: repo, validator: validator.New()}
}

var ErrPlanAlreadyExists = errors.New("plan already exists")

func (s *planService) CreatePlan(ctx context.Context, input dto.CreatePlanRequest) error {
	if err := s.validator.Struct(input); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	existingPlan, err := s.repo.GetPlanByName(ctx, input.Name)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if err == nil && existingPlan.Name != "" {
		return ErrPlanAlreadyExists
	}

	plan := db.CreatePlanParams{
		Name:                  input.Name,
		PriceMonthly:          int32(existingPlan.PriceMonthly),
		MaxWebsites:           input.MaxWebsites,
		CheckIntervalSeconds:  input.CheckIntervalSeconds,
		HasPerformanceReports: input.HasPerformanceReports,
		HasSeoAudits:          input.HasSEOAudits,
		HasPublicStatusPage:   input.HasPublicStatusPage,
	}

	if err := s.repo.CreatePlan(ctx, plan); err != nil {
		return err
	}

	return nil
}
