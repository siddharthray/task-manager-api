package model

import "time"

type Task struct {
	ID          int64      `json:"id"`
	Text        string     `json:"text"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ReopenedAt  *time.Time `json:"reopened_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UserID      *int64     `json:"user_id,omitempty"`
}
