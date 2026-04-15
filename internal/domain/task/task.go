package task

import (
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	RepeatType     string     `json:"repeat_type,omitempty"`
	RepeatConfig   []byte     `json:"repeat_config,omitempty"`
	NextOccurrence *time.Time `json:"next_occurrence,omitempty"`
	ParentID       *int64     `json:"parent_id,omitempty"`
	RepeatTime     string     `json:"repeat_time,omitempty"`

	RepeatUntil *time.Time `json:"repeat_until,omitempty"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
