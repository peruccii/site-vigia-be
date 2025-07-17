package models

type Plan struct {
	ID                    int     `json:"id" db:"id"`
	Name                  string  `json:"name" db:"name"` // free, freelancer ,agency
	PriceMonthly          float64 `json:"price_monthly" db:"price_monthly"`
	MaxWebsites           int     `json:"max_websites" db:"max_websites"`
	CheckIntervalSeconds  int     `json:"check_interval_seconds" db:"check_interval_seconds"`
	HasPerformanceReports bool    `json:"has_performance_reports" db:"has_performance_reports"`
	HasSEOAudits          bool    `json:"has_seo_audits" db:"has_seo_audits"`
	HasPublicStatusPage   bool    `json:"has_public_status_page" db:"has_public_status_page"`
}
