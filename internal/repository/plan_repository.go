package repository

import (
	"context"

	"peruccii/site-vigia-be/db"
)

type PlanRepository struct {
	queries *db.Queries
}

func NewPlanRepository(queries *db.Queries) *PlanRepository {
	return &PlanRepository{
		queries: queries,
	}
}

func (r *PlanRepository) CreatePlan(ctx context.Context, input db.CreatePlanParams) error {
	return r.queries.CreatePlan(ctx, input)
}

func (r *PlanRepository) GetPlanByName(ctx context.Context, name string) (db.Plan, error) {
	return r.queries.GetPlanByName(ctx, name)
}
