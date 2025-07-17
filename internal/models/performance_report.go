package models

import (
	"time"
	"github.com/google/uuid"
)

type PerformanceReport struct {
	ID             int32     `json:"id" db:"id"`
	WebsiteID      uuid.UUID `json:"website_id" db:"website_id"`
	CheckedAt      time.Time `json:"checked_at" db:"checked_at"`
	TTFBMs         *int      `json:"ttfb_ms" db:"ttfb_ms"`
	LCPMs          *int      `json:"lcp_ms" db:"lcp_ms"`
	FullLoadTimeMs *int      `json:"full_load_time_ms" db:"full_load_time_ms"`
	PageSizeKB     *int      `json:"page_size_kb" db:"page_size_kb"`
}
