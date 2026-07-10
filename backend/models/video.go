package models

import "time"

type VideoStatus string

const (
	StatusPending    VideoStatus = "pending"
	StatusProcessing VideoStatus = "processing"
	StatusCompleted  VideoStatus = "completed"
	StatusFailed     VideoStatus = "failed"
)

type Video struct {
	ID               string      `json:"id" db:"id"`
	UserID           string      `json:"user_id" db:"user_id"`
	Title            *string     `json:"title" db:"title"`
	OriginalFilename *string     `json:"original_filename" db:"original_filename"`
	StorageKey       *string     `json:"-" db:"storage_key"`
	YouTubeURL       *string     `json:"youtube_url,omitempty" db:"youtube_url"`
	DurationSeconds  *int        `json:"duration_seconds" db:"duration_seconds"`
	Language         string      `json:"language" db:"language"`
	FileSizeBytes    *int64      `json:"file_size_bytes" db:"file_size_bytes"`
	Status           VideoStatus `json:"status" db:"status"`
	ErrorMessage     *string     `json:"error_message,omitempty" db:"error_message"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	CompletedAt      *time.Time  `json:"completed_at,omitempty" db:"completed_at"`
}

type UploadVideoRequest struct {
	YouTubeURL string `json:"youtube_url"`
	Language   string `json:"language" validate:"oneof=en ru kk auto"`
}

type VideoListResponse struct {
	Videos     []Video `json:"videos"`
	TotalCount int     `json:"total_count"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
}

// VideoDetail is returned by GET /videos/:id — includes summary content and QA.
type VideoDetail struct {
	ID               string          `json:"id"`
	UserID           string          `json:"user_id"`
	Title            *string         `json:"title"`
	OriginalFilename *string         `json:"original_filename,omitempty"`
	YouTubeURL       *string         `json:"youtube_url,omitempty"`
	DurationSeconds  *int            `json:"duration_seconds"`
	Language         string          `json:"language"`
	FileSizeBytes    *int64          `json:"file_size_bytes,omitempty"`
	Status           VideoStatus     `json:"status"`
	ErrorMessage     *string         `json:"error_message,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	CompletedAt      *time.Time      `json:"completed_at,omitempty"`
	Summary          *SummaryContent `json:"summary,omitempty"`
	QA               []QAItem        `json:"qa,omitempty"`
}
