package handlers

import (
	"encoding/json"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Status       taskdomain.Status `json:"status"`
	RepeatType   string            `json:"repeat_type,omitempty"`
	RepeatConfig json.RawMessage   `json:"repeat_config,omitempty"`
	RepeatUntil  *time.Time        `json:"repeat_until,omitempty"`
	RepeatTime   string            `json:"repeat_time,omitempty"`
}

type taskDTO struct {
	ID             int64             `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Status         taskdomain.Status `json:"status"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	RepeatType     string            `json:"repeat_type,omitempty"`
	RepeatConfig   json.RawMessage   `json:"repeat_config,omitempty"`
	NextOccurrence *time.Time        `json:"next_occurrence,omitempty"`
	RepeatUntil    *time.Time        `json:"repeat_until,omitempty"`

	RepeatTime string `json:"repeat_time,omitempty"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:             task.ID,
		Title:          task.Title,
		Description:    task.Description,
		Status:         task.Status,
		CreatedAt:      task.CreatedAt,
		UpdatedAt:      task.UpdatedAt,
		RepeatType:     task.RepeatType,
		RepeatConfig:   task.RepeatConfig,
		NextOccurrence: task.NextOccurrence,
		RepeatUntil:    task.RepeatUntil,
		RepeatTime:     task.RepeatTime,
	}
}
