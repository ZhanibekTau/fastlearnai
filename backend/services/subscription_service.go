package services

import (
	"context"
	"net/http"

	"github.com/fastlearn/backend/config"
	"github.com/fastlearn/backend/models"
	"github.com/fastlearn/backend/repository"
	"github.com/stripe/stripe-go/v76/webhook"
)

type SubscriptionService struct {
	subs repository.SubscriptionRepository
	cfg  *config.Config
}

func NewSubscriptionService(subs repository.SubscriptionRepository, cfg *config.Config) *SubscriptionService {
	return &SubscriptionService{subs: subs, cfg: cfg}
}

func (s *SubscriptionService) GetMine(ctx context.Context, userID string) (*models.Subscription, error) {
	return s.subs.FindByUserID(ctx, userID)
}

func (s *SubscriptionService) GetUsage(ctx context.Context, userID string) (map[string]any, error) {
	plan, err := s.subs.GetPlanByUserID(ctx, userID)
	if err != nil {
		plan = models.PlanFree
	}

	usage, err := s.subs.GetUsage(ctx, userID, repository.CurrentMonth())
	if err != nil {
		return nil, err
	}

	limits := models.PlanConfig[plan]
	return map[string]any{
		"plan":          string(plan),
		"videos_count":  usage.VideosCount,
		"total_minutes": usage.TotalMinutes,
		"limits": map[string]any{
			"max_videos":        limits.MaxVideosPerMonth,
			"max_duration_min":  limits.MaxDurationMin,
			"has_pdf":           limits.HasPDF,
			"has_markdown":      limits.HasMarkdown,
			"has_qa":            limits.HasQA,
		},
	}, nil
}

func (s *SubscriptionService) HandleStripeWebhook(r *http.Request) error {
	if s.cfg.StripeWebhookSecret == "" {
		return nil
	}
	payload := make([]byte, r.ContentLength)
	_, _ = r.Body.Read(payload)
	_, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), s.cfg.StripeWebhookSecret)
	return err
}
