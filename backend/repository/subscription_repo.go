package repository

import (
	"context"
	"errors"
	"time"

	"github.com/fastlearn/backend/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository interface {
	FindByUserID(ctx context.Context, userID string) (*models.Subscription, error)
	CreateDefault(ctx context.Context, userID string) error
	GetPlanByUserID(ctx context.Context, userID string) (models.Plan, error)
	GetUsage(ctx context.Context, userID, month string) (*models.Usage, error)
	IncrementUsage(ctx context.Context, userID, month string) error
	AdminSetPlan(ctx context.Context, userID string, plan models.Plan, periodStart, periodEnd *time.Time) (*models.Subscription, error)
}

type pgSubscriptionRepo struct{ db *pgxpool.Pool }

func NewSubscriptionRepo(db *pgxpool.Pool) SubscriptionRepository {
	return &pgSubscriptionRepo{db: db}
}

func (r *pgSubscriptionRepo) FindByUserID(ctx context.Context, userID string) (*models.Subscription, error) {
	var s models.Subscription
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, plan, stripe_subscription_id, status,
		        current_period_start, current_period_end, created_at
		 FROM subscriptions WHERE user_id=$1`, userID,
	).Scan(&s.ID, &s.UserID, &s.Plan, &s.StripeSubscriptionID, &s.Status,
		&s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return &models.Subscription{UserID: userID, Plan: models.PlanFree, Status: "active"}, nil
	}
	return &s, err
}

func (r *pgSubscriptionRepo) CreateDefault(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO subscriptions (id, user_id, plan, status) VALUES ($1,$2,'free','active')
		 ON CONFLICT (user_id) DO NOTHING`,
		uuid.New().String(), userID,
	)
	return err
}

// AdminSetPlan upserts the user's subscription to the given plan/period, for manual grants by an admin.
func (r *pgSubscriptionRepo) AdminSetPlan(ctx context.Context, userID string, plan models.Plan, periodStart, periodEnd *time.Time) (*models.Subscription, error) {
	var s models.Subscription
	err := r.db.QueryRow(ctx,
		`INSERT INTO subscriptions (id, user_id, plan, status, current_period_start, current_period_end)
		 VALUES ($1,$2,$3,'active',$4,$5)
		 ON CONFLICT (user_id) DO UPDATE
		   SET plan=$3, status='active', current_period_start=$4, current_period_end=$5
		 RETURNING id, user_id, plan, stripe_subscription_id, status, current_period_start, current_period_end, created_at`,
		uuid.New().String(), userID, string(plan), periodStart, periodEnd,
	).Scan(&s.ID, &s.UserID, &s.Plan, &s.StripeSubscriptionID, &s.Status,
		&s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *pgSubscriptionRepo) GetPlanByUserID(ctx context.Context, userID string) (models.Plan, error) {
	var plan string
	err := r.db.QueryRow(ctx,
		`SELECT plan FROM subscriptions WHERE user_id=$1 AND status='active'`, userID,
	).Scan(&plan)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.PlanFree, nil
	}
	return models.Plan(plan), err
}

func (r *pgSubscriptionRepo) GetUsage(ctx context.Context, userID, month string) (*models.Usage, error) {
	var u models.Usage
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, month, videos_count, total_minutes FROM usage WHERE user_id=$1 AND month=$2`,
		userID, month,
	).Scan(&u.ID, &u.UserID, &u.Month, &u.VideosCount, &u.TotalMinutes)
	if errors.Is(err, pgx.ErrNoRows) {
		return &models.Usage{UserID: userID}, nil
	}
	return &u, err
}

func (r *pgSubscriptionRepo) IncrementUsage(ctx context.Context, userID, month string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO usage (id, user_id, month, videos_count, total_minutes)
		 VALUES ($1,$2,$3,1,0)
		 ON CONFLICT (user_id, month) DO UPDATE SET videos_count = usage.videos_count + 1`,
		uuid.New().String(), userID, month,
	)
	return err
}

func CurrentMonth() string {
	return time.Now().Format("2006-01-01")
}
