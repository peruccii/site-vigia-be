package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)
// JSONB
type SEOResults map[string]any // string key and any value

func (s SEOResults) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *SEOResults) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return errors.New("cannot scan non-string value into SEOResults")
	}
}

type SEOAudit struct {
	ID        int32      `json:"id" db:"id"`
	WebsiteID uuid.UUID  `json:"website_id" db:"website_id"`
	AuditedAt time.Time  `json:"audited_at" db:"audited_at"`
	Results   SEOResults `json:"results" db:"results"`
}
