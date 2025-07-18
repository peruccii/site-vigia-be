package repository

import (
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
