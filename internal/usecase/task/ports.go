package task

import (
	"context"
	"encoding/json"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
	GetDueRecurringTasks(ctx context.Context, now time.Time) ([]taskdomain.Task, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
	GenerateNextTasks(ctx context.Context) error
}

type CreateInput struct {
	Title        string
	Description  string
	Status       taskdomain.Status
	RepeatType   string          `json:"repeat_type,omitempty"`
	RepeatConfig json.RawMessage `json:"repeat_config,omitempty"`
	RepeatUntil  *time.Time      `json:"repeat_until,omitempty"`
	RepeatTime   string          `json:"repeat_time,omitempty"`
}

type UpdateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
}
