package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/fastlearn/backend/config"
	"github.com/fastlearn/backend/models"
	"github.com/fastlearn/backend/repository"
	"github.com/redis/go-redis/v9"
)

func mergeQA(content, qa json.RawMessage) (json.RawMessage, error) {
	if len(qa) == 0 || string(qa) == "null" {
		return content, nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(content, &m); err != nil {
		return nil, err
	}
	m["qa"] = qa
	return json.Marshal(m)
}

type Processor struct {
	videos     repository.VideoRepository
	summaries  repository.SummaryRepository
	processing repository.ProcessingRepository
	subs       repository.SubscriptionRepository
	rdb        *redis.Client
	cfg        *config.Config
	workers    int
}

func NewProcessor(
	videos repository.VideoRepository,
	summaries repository.SummaryRepository,
	processing repository.ProcessingRepository,
	subs repository.SubscriptionRepository,
	rdb *redis.Client,
	cfg *config.Config,
) *Processor {
	return &Processor{
		videos:     videos,
		summaries:  summaries,
		processing: processing,
		subs:       subs,
		rdb:        rdb,
		cfg:        cfg,
		workers:    4,
	}
}

func (p *Processor) Start() {
	slog.Info("queue processor ready", "workers", p.workers)
	for i := 0; i < p.workers; i++ {
		go p.runWorker(i)
	}
}

func (p *Processor) runWorker(id int) {
	for {
		job, err := Pop(context.Background(), p.rdb)
		if err != nil {
			slog.Error("queue pop failed", "worker", id, "error", err)
			time.Sleep(5 * time.Second)
			continue
		}
		p.process(job)
	}
}

func (p *Processor) process(job *Job) {
	ctx := context.Background()
	start := time.Now()

	src := "file"
	if job.YouTubeURL != "" {
		src = "youtube"
	}
	slog.Info("processing started", "video", job.VideoID, "src", src, "lang", job.Language)

	if err := p.videos.UpdateStatus(ctx, job.VideoID, models.StatusProcessing, ""); err != nil {
		slog.Error("status update failed", "video", job.VideoID, "error", err)
		return
	}

	result, err := p.processing.Process(ctx, repository.ProcessRequest{
		VideoID:    job.VideoID,
		FilePath:   job.FilePath,
		YouTubeURL: job.YouTubeURL,
		Language:   job.Language,
	})
	if err != nil {
		slog.Error("processing failed", "video", job.VideoID, "src", src, "error", err)
		_ = p.videos.UpdateStatus(ctx, job.VideoID, models.StatusFailed, err.Error())
		return
	}

	if exceeded, limitMin := p.durationExceedsLimit(ctx, job.UserID, result.DurationSeconds); exceeded {
		msg := fmt.Sprintf("video duration exceeds the %d-minute limit for your plan", limitMin)
		slog.Warn("duration limit exceeded", "video", job.VideoID, "duration_sec", result.DurationSeconds, "limit_min", limitMin)
		_ = p.videos.UpdateStatus(ctx, job.VideoID, models.StatusFailed, msg)
		if job.FilePath != "" {
			_ = os.Remove(job.FilePath)
		}
		return
	}

	merged, err := mergeQA(result.Content, result.QA)
	if err != nil {
		slog.Warn("qa merge failed, storing without qa", "error", err)
		merged = result.Content
	}
	if err := p.summaries.Upsert(ctx, job.VideoID, merged, result.Transcript); err != nil {
		slog.Error("summary store failed", "video", job.VideoID, "error", err)
		_ = p.videos.UpdateStatus(ctx, job.VideoID, models.StatusFailed, "failed to store summary")
		return
	}

	if err := p.videos.UpdateCompleted(ctx, job.VideoID, result.Title, result.DurationSeconds); err != nil {
		slog.Error("completion update failed", "video", job.VideoID, "error", err)
		return
	}

	if job.FilePath != "" {
		if err := os.Remove(job.FilePath); err != nil && !os.IsNotExist(err) {
			slog.Warn("file cleanup failed", "path", job.FilePath, "error", err)
		}
	}

	elapsed := time.Since(start).Round(time.Second)
	slog.Info("processing done", "video", job.VideoID, "title", result.Title, "elapsed", elapsed)
}

// durationExceedsLimit checks the video's actual duration against the plan's MaxDurationMin.
// Duration is only known after processing (YouTube/Whisper), so this runs post-hoc rather than at upload time.
func (p *Processor) durationExceedsLimit(ctx context.Context, userID string, durationSeconds int) (bool, int) {
	plan, err := p.subs.GetPlanByUserID(ctx, userID)
	if err != nil {
		plan = models.PlanFree
	}
	limitMin := models.PlanConfig[plan].MaxDurationMin
	if limitMin == -1 {
		return false, 0
	}
	if durationSeconds > limitMin*60 {
		return true, limitMin
	}
	return false, 0
}
