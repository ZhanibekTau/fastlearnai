package models

import "time"

type KeyPoint struct {
	Topic       string `json:"topic"`
	Explanation string `json:"explanation"`
}

type Example struct {
	Topic   string `json:"topic"`
	Example string `json:"example"`
}

type QAItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type SummaryContent struct {
	Overview     string     `json:"overview"`
	KeyPoints    []KeyPoint `json:"key_points"`
	KeyTakeaways []string   `json:"key_takeaways"`
	Examples     []Example  `json:"examples"`
}

type Summary struct {
	ID         string          `json:"id" db:"id"`
	VideoID    string          `json:"video_id" db:"video_id"`
	Title      string          `json:"title"`
	Duration   string          `json:"duration"`
	Language   string          `json:"language"`
	Content    *SummaryContent `json:"summary" db:"content"`
	QA         []QAItem        `json:"qa,omitempty"`
	Transcript *string         `json:"transcript,omitempty" db:"transcript"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}
