package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fastlearn/backend/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SummaryRepository interface {
	Upsert(ctx context.Context, videoID string, content json.RawMessage, transcript string) error
	FindByVideoID(ctx context.Context, videoID string) (*models.Summary, error)
}

type pgSummaryRepo struct{ db *pgxpool.Pool }

func NewSummaryRepo(db *pgxpool.Pool) SummaryRepository {
	return &pgSummaryRepo{db: db}
}

func (r *pgSummaryRepo) Upsert(ctx context.Context, videoID string, content json.RawMessage, transcript string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO summaries (id, video_id, content, transcript)
		 VALUES ($1,$2,$3,$4)
		 ON CONFLICT (video_id) DO UPDATE SET content=$3, transcript=$4`,
		uuid.New().String(), videoID, content, transcript,
	)
	return err
}

func (r *pgSummaryRepo) FindByVideoID(ctx context.Context, videoID string) (*models.Summary, error) {
	var s models.Summary
	var raw json.RawMessage
	err := r.db.QueryRow(ctx,
		`SELECT id, video_id, content, transcript, created_at FROM summaries WHERE video_id=$1`, videoID,
	).Scan(&s.ID, &s.VideoID, &raw, &s.Transcript, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if len(raw) > 0 {
		var payload struct {
			models.SummaryContent
			QA []models.QAItem `json:"qa"`
		}
		if err := json.Unmarshal(raw, &payload); err == nil {
			s.Content = &payload.SummaryContent
			s.QA = payload.QA
		}
	}
	return &s, nil
}
