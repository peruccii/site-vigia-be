package models

import (
	"time"
	"github.com/google/uuid"
)

type Incident struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	WebsiteID       uuid.UUID  `json:"website_id" db:"website_id"`
	StartedAt       time.Time  `json:"started_at" db:"started_at"`
	EndedAt         *time.Time `json:"ended_at" db:"ended_at"`
	DurationSeconds *int       `json:"duration_seconds" db:"duration_seconds"`
	Cause           string     `json:"cause" db:"cause"`
}
