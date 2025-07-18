package repository

import (
	"context"
	"database/sql"

	"peruccii/site-vigia-be/db"
)

type PlanRepository struct {
	queries *db.Queries
}

func NewPlanRepository(dbConn *sql.DB) *PlanRepository {
	return &PlanRepository{
		queries: db.New(dbConn),
	}
}

func (r *PlanRepository) CreatePlan(ctx context.Context, input db.CreatePlanParams) error {
	return r.queries.CreatePlan(ctx, input)
}

func (r *PlanRepository) GetPlanByName(ctx context.Context, name string) (db.Plan, error) {
	return r.queries.GetPlanByName(ctx, name)
}
