package dto

type CreatePlanRequest struct {
	Name                  string  `json:"name" validate:"required,min=1,max=100"`
	PriceMonthly          float64 `json:"price_monthly" validate:"required,min=0.01"`
	StripePriceID         string  `json:"stripe_price_id" validate:"required"`
	MaxWebsites           int     `json:"max_websites" validate:"required,min=1,max=100"`
	CheckIntervalSeconds  int     `json:"check_interval_seconds" validate:"required,min=1,max=100"`
	HasPerformanceReports bool    `json:"has_performance_reports" validate:"required"`
	HasSEOAudits          bool    `json:"has_seo_audits" validate:"required"`
	HasPublicStatusPage   bool    `json:"has_public_status_page" validate:"required"`
}
