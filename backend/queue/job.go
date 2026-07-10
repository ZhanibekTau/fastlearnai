package queue

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

const QueueKey = "fastlearn:jobs"

type Job struct {
	VideoID    string `json:"video_id"`
	UserID     string `json:"user_id"`
	Language   string `json:"language"`
	FilePath   string `json:"file_path,omitempty"`   // local /uploads/video path
	YouTubeURL string `json:"youtube_url,omitempty"`
	Priority   int    `json:"priority"`
}

func Push(ctx context.Context, rdb *redis.Client, job Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return rdb.LPush(ctx, QueueKey, data).Err()
}

func Pop(ctx context.Context, rdb *redis.Client) (*Job, error) {
	result, err := rdb.BRPop(ctx, 0, QueueKey).Result()
	if err != nil {
		return nil, err
	}
	var job Job
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, err
	}
	return &job, nil
}
