package models

import (
	"time"
	"github.com/google/uuid"
)

type UptimeCheck struct {
	ID             int32     `json:"id" db:"id"`
	WebsiteID      uuid.UUID `json:"website_id" db:"website_id"`
	CheckedAt      time.Time `json:"checked_at" db:"checked_at"`
	IsUp           bool      `json:"is_up" db:"is_up"`
	ResponseTimeMs int       `json:"response_time_ms" db:"response_time_ms"`
	StatusCode     *int      `json:"status_code" db:"status_code"`
	ErrorMessage   *string   `json:"error_message" db:"error_message"`
}
