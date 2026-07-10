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

type CreateVideoInput struct {
	UserID           string
	Title            *string
	OriginalFilename *string
	StorageKey       *string
	YouTubeURL       *string
	Language         string
	FileSizeBytes    *int64
}

type VideoRepository interface {
	Create(ctx context.Context, input CreateVideoInput) (*models.Video, error)
	FindByID(ctx context.Context, id string) (*models.Video, error)
	FindByIDAndUserID(ctx context.Context, id, userID string) (*models.Video, error)
	ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]models.Video, int, error)
	UpdateStorageKey(ctx context.Context, id, path string) error
	UpdateStatus(ctx context.Context, id string, status models.VideoStatus, errMsg string) error
	UpdateCompleted(ctx context.Context, id, title string, durationSeconds int) error
	Delete(ctx context.Context, id, userID string) error
}

type pgVideoRepo struct{ db *pgxpool.Pool }

func NewVideoRepo(db *pgxpool.Pool) VideoRepository {
	return &pgVideoRepo{db: db}
}

func (r *pgVideoRepo) Create(ctx context.Context, input CreateVideoInput) (*models.Video, error) {
	id := uuid.New().String()
	now := time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO videos
		 (id, user_id, title, original_filename, storage_key, youtube_url, language, file_size_bytes, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,'pending',$9)`,
		id, input.UserID, input.Title, input.OriginalFilename,
		input.StorageKey, input.YouTubeURL, input.Language, input.FileSizeBytes, now,
	)
	if err != nil {
		return nil, err
	}

	return &models.Video{
		ID:               id,
		UserID:           input.UserID,
		Title:            input.Title,
		OriginalFilename: input.OriginalFilename,
		StorageKey:       input.StorageKey,
		YouTubeURL:       input.YouTubeURL,
		Language:         input.Language,
		FileSizeBytes:    input.FileSizeBytes,
		Status:           models.StatusPending,
		CreatedAt:        now,
	}, nil
}

func (r *pgVideoRepo) FindByID(ctx context.Context, id string) (*models.Video, error) {
	return r.scanVideo(r.db.QueryRow(ctx,
		`SELECT id, user_id, title, original_filename, storage_key, youtube_url,
		        duration_seconds, language, file_size_bytes, status, error_message, created_at, completed_at
		 FROM videos WHERE id=$1`, id,
	))
}

func (r *pgVideoRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*models.Video, error) {
	return r.scanVideo(r.db.QueryRow(ctx,
		`SELECT id, user_id, title, original_filename, storage_key, youtube_url,
		        duration_seconds, language, file_size_bytes, status, error_message, created_at, completed_at
		 FROM videos WHERE id=$1 AND user_id=$2`, id, userID,
	))
}

func (r *pgVideoRepo) ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]models.Video, int, error) {
	offset := (page - 1) * pageSize
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, title, original_filename, storage_key, youtube_url,
		        duration_seconds, language, file_size_bytes, status, error_message, created_at, completed_at
		 FROM videos WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		v, err := r.scanVideo(rows)
		if err != nil {
			continue
		}
		videos = append(videos, *v)
	}

	var total int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM videos WHERE user_id=$1`, userID).Scan(&total)

	return videos, total, nil
}

func (r *pgVideoRepo) UpdateStorageKey(ctx context.Context, id, path string) error {
	_, err := r.db.Exec(ctx, `UPDATE videos SET storage_key=$1 WHERE id=$2`, path, id)
	return err
}

func (r *pgVideoRepo) UpdateStatus(ctx context.Context, id string, status models.VideoStatus, errMsg string) error {
	var q string
	var args []any
	if errMsg != "" {
		q = `UPDATE videos SET status=$1, error_message=$2 WHERE id=$3`
		args = []any{status, errMsg, id}
	} else {
		q = `UPDATE videos SET status=$1, error_message=NULL WHERE id=$2`
		args = []any{status, id}
	}
	_, err := r.db.Exec(ctx, q, args...)
	return err
}

func (r *pgVideoRepo) UpdateCompleted(ctx context.Context, id, title string, durationSeconds int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE videos SET status='completed', title=$1, duration_seconds=$2, completed_at=$3, error_message=NULL WHERE id=$4`,
		title, durationSeconds, time.Now(), id,
	)
	return err
}

func (r *pgVideoRepo) Delete(ctx context.Context, id, userID string) error {
	res, err := r.db.Exec(ctx,
		`DELETE FROM videos WHERE id=$1 AND user_id=$2`, id, userID,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// scanVideo works with both pgx.Row and pgx.Rows via the Scanner interface.
type scanner interface {
	Scan(dest ...any) error
}

func (r *pgVideoRepo) scanVideo(row scanner) (*models.Video, error) {
	var v models.Video
	err := row.Scan(
		&v.ID, &v.UserID, &v.Title, &v.OriginalFilename, &v.StorageKey, &v.YouTubeURL,
		&v.DurationSeconds, &v.Language, &v.FileSizeBytes, &v.Status, &v.ErrorMessage,
		&v.CreatedAt, &v.CompletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &v, err
}
