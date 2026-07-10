package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/fastlearn/backend/config"
	"github.com/fastlearn/backend/models"
	"github.com/fastlearn/backend/queue"
	"github.com/fastlearn/backend/repository"
	"github.com/redis/go-redis/v9"
)

var (
	ErrLimitReached = errors.New("monthly limit reached")
	ErrNotFound     = errors.New("not found")
	ErrForbidden    = errors.New("forbidden")
)

type UploadInput struct {
	UserID     string
	File       multipart.File
	FileHeader *multipart.FileHeader
	YouTubeURL string
	Language   string
}

type VideoService struct {
	videos    repository.VideoRepository
	summaries repository.SummaryRepository
	subs      repository.SubscriptionRepository
	rdb       *redis.Client
	cfg       *config.Config
}

func NewVideoService(
	videos repository.VideoRepository,
	summaries repository.SummaryRepository,
	subs repository.SubscriptionRepository,
	rdb *redis.Client,
	cfg *config.Config,
) *VideoService {
	return &VideoService{videos: videos, summaries: summaries, subs: subs, rdb: rdb, cfg: cfg}
}

func (s *VideoService) Upload(ctx context.Context, input UploadInput) (*models.Video, error) {
	if err := s.enforceLimit(ctx, input.UserID); err != nil {
		return nil, err
	}

	in := repository.CreateVideoInput{
		UserID:   input.UserID,
		Language: input.Language,
	}

	var filePath string

	if input.YouTubeURL != "" {
		in.YouTubeURL = &input.YouTubeURL
	} else if input.File != nil && input.FileHeader != nil {
		name := input.FileHeader.Filename
		size := input.FileHeader.Size
		in.OriginalFilename = &name
		in.FileSizeBytes = &size
	}

	// Создаём запись в БД чтобы получить video.ID
	video, err := s.videos.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	// Сохраняем файл на диск под именем video.ID
	if input.File != nil && input.FileHeader != nil {
		ext := filepath.Ext(input.FileHeader.Filename)
		filePath, err = s.saveFile(input.File, video.ID, ext)
		if err != nil {
			return nil, fmt.Errorf("save file: %w", err)
		}
		// Обновляем storage_key в БД (локальный путь)
		if err := s.videos.UpdateStorageKey(ctx, video.ID, filePath); err != nil {
			return nil, err
		}
		video.StorageKey = &filePath
	}

	_ = s.subs.IncrementUsage(ctx, input.UserID, repository.CurrentMonth())

	job := queue.Job{
		VideoID:    video.ID,
		UserID:     input.UserID,
		Language:   input.Language,
		FilePath:   filePath,
		YouTubeURL: input.YouTubeURL,
	}
	_ = queue.Push(ctx, s.rdb, job)

	return video, nil
}

// saveFile записывает multipart файл в UPLOADS_DIR/{videoID}{ext} и возвращает путь.
func (s *VideoService) saveFile(file multipart.File, videoID, ext string) (string, error) {
	if err := os.MkdirAll(s.cfg.UploadsDir, 0755); err != nil {
		return "", err
	}
	dest := filepath.Join(s.cfg.UploadsDir, videoID+ext)
	out, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}
	return dest, nil
}

func (s *VideoService) List(ctx context.Context, userID string, page, pageSize int) (*models.VideoListResponse, error) {
	if page < 1 {
		page = 1
	}
	videos, total, err := s.videos.ListByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &models.VideoListResponse{
		Videos:     videos,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *VideoService) Get(ctx context.Context, id, userID string) (*models.VideoDetail, error) {
	v, err := s.videos.FindByIDAndUserID(ctx, id, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return s.toDetail(ctx, v), nil
}

func (s *VideoService) toDetail(ctx context.Context, v *models.Video) *models.VideoDetail {
	d := &models.VideoDetail{
		ID:               v.ID,
		UserID:           v.UserID,
		Title:            v.Title,
		OriginalFilename: v.OriginalFilename,
		YouTubeURL:       v.YouTubeURL,
		DurationSeconds:  v.DurationSeconds,
		Language:         v.Language,
		FileSizeBytes:    v.FileSizeBytes,
		Status:           v.Status,
		ErrorMessage:     v.ErrorMessage,
		CreatedAt:        v.CreatedAt,
		CompletedAt:      v.CompletedAt,
	}
	if v.Status == models.StatusCompleted {
		if sum, err := s.summaries.FindByVideoID(ctx, v.ID); err == nil {
			d.Summary = sum.Content
			d.QA = sum.QA
		}
	}
	return d
}

func (s *VideoService) Delete(ctx context.Context, id, userID string) error {
	return s.videos.Delete(ctx, id, userID)
}

func (s *VideoService) Regenerate(ctx context.Context, id, userID string) error {
	v, err := s.videos.FindByIDAndUserID(ctx, id, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	if err := s.videos.UpdateStatus(ctx, id, models.StatusPending, ""); err != nil {
		return err
	}

	filePath := ""
	if v.StorageKey != nil {
		filePath = *v.StorageKey
	}

	return queue.Push(ctx, s.rdb, queue.Job{
		VideoID:    id,
		UserID:     userID,
		Language:   v.Language,
		FilePath:   filePath,
		YouTubeURL: func() string {
			if v.YouTubeURL != nil {
				return *v.YouTubeURL
			}
			return ""
		}(),
	})
}

func (s *VideoService) GetPublic(ctx context.Context, id string) (*models.VideoDetail, error) {
	v, err := s.videos.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if v.Status != models.StatusCompleted {
		return nil, ErrNotFound
	}
	return s.toDetail(ctx, v), nil
}

func (s *VideoService) enforceLimit(ctx context.Context, userID string) error {
	plan, err := s.subs.GetPlanByUserID(ctx, userID)
	if err != nil {
		plan = models.PlanFree
	}
	limits := models.PlanConfig[plan]
	if limits.MaxVideosPerMonth == -1 {
		return nil
	}
	usage, err := s.subs.GetUsage(ctx, userID, repository.CurrentMonth())
	if err != nil {
		return err
	}
	if usage.VideosCount >= limits.MaxVideosPerMonth {
		return ErrLimitReached
	}
	return nil
}
