package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ProcessRequest is sent to the FastAPI AI service.
type ProcessRequest struct {
	VideoID    string `json:"video_id"`
	FilePath   string `json:"file_path,omitempty"`   // local /uploads/video path
	YouTubeURL string `json:"youtube_url,omitempty"`
	Language   string `json:"language"`
}

// ProcessResult is returned by the FastAPI AI service.
type ProcessResult struct {
	Title           string          `json:"title"`
	DurationSeconds int             `json:"duration_seconds"`
	Language        string          `json:"language"`
	Content         json.RawMessage `json:"content"`
	QA              json.RawMessage `json:"qa"`
	Transcript      string          `json:"transcript"`
}

type ProcessingRepository interface {
	Process(ctx context.Context, req ProcessRequest) (*ProcessResult, error)
}

type httpProcessingRepo struct {
	baseURL string
	client  *http.Client
}

func NewProcessingRepo(fastAPIURL string) ProcessingRepository {
	return &httpProcessingRepo{
		baseURL: fastAPIURL,
		client:  &http.Client{Timeout: 30 * time.Minute},
	}
}

func (r *httpProcessingRepo) Process(ctx context.Context, req ProcessRequest) (*ProcessResult, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/process", bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := r.client.Do(httpReq)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt*attempt) * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("processing service returned %d", resp.StatusCode)
			time.Sleep(time.Duration(attempt*attempt) * time.Second)
			continue
		}

		var result ProcessResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, lastErr
}
