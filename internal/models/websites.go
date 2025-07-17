package models

import (
	"time"

	"github.com/google/uuid"
)

type Website struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	UserID               uuid.UUID  `json:"user_id" db:"user_id"`
	Name                 string     `json:"name" db:"name"`
	URL                  string     `json:"url" db:"url"`
	IsActive             bool       `json:"is_active" db:"is_active"`
	CheckIntervalSeconds int        `json:"check_interval_seconds" db:"check_interval_seconds"`
	LastCheckedAt        *time.Time `json:"last_checked_at" db:"last_checked_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

type WebsiteWithUser struct {
	Website
	UserName  string `json:"user_name" db:"user_name"`
	UserEmail string `json:"user_email" db:"user_email"`
}
