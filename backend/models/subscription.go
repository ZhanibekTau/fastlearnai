package models

import "time"

type Plan string

const (
	PlanFree    Plan = "free"
	PlanPro     Plan = "pro"
	PlanPremium Plan = "premium"
)

type Subscription struct {
	ID                   string    `json:"id" db:"id"`
	UserID               string    `json:"user_id" db:"user_id"`
	Plan                 Plan      `json:"plan" db:"plan"`
	StripeSubscriptionID *string   `json:"stripe_subscription_id,omitempty" db:"stripe_subscription_id"`
	Status               string    `json:"status" db:"status"`
	CurrentPeriodStart   *time.Time `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end" db:"current_period_end"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
}

type Usage struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Month       time.Time `json:"month" db:"month"`
	VideosCount int       `json:"videos_count" db:"videos_count"`
	TotalMinutes int      `json:"total_minutes" db:"total_minutes"`
}

type PlanLimits struct {
	MaxVideosPerMonth int
	MaxDurationMin    int
	HasPDF            bool
	HasMarkdown       bool
	HasQA             bool
	HasTranscript     bool
	HasAPI            bool
	HistoryDays       int
}

var PlanConfig = map[Plan]PlanLimits{
	PlanFree: {
		MaxVideosPerMonth: 3,
		MaxDurationMin:    180,
		HasPDF:            false,
		HasMarkdown:       false,
		HasQA:             false,
		HasTranscript:     false,
		HasAPI:            false,
		HistoryDays:       7,
	},
	PlanPro: {
		MaxVideosPerMonth: 30,
		MaxDurationMin:    180,
		HasPDF:            true,
		HasMarkdown:       true,
		HasQA:             true,
		HasTranscript:     true,
		HasAPI:            false,
		HistoryDays:       30,
	},
	PlanPremium: {
		MaxVideosPerMonth: -1, // unlimited
		MaxDurationMin:    180,
		HasPDF:            true,
		HasMarkdown:       true,
		HasQA:             true,
		HasTranscript:     true,
		HasAPI:            true,
		HistoryDays:       -1,
	},
}
