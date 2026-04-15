package task

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:        normalized.Title,
		Description:  normalized.Description,
		Status:       normalized.Status,
		RepeatType:   input.RepeatType,
		RepeatConfig: input.RepeatConfig,
		RepeatUntil:  input.RepeatUntil,
		RepeatTime:   input.RepeatTime,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		UpdatedAt:   s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func (s *Service) GenerateNextTasks(ctx context.Context) error {
	now := s.now()

	tasks, err := s.repo.GetDueRecurringTasks(ctx, now)
	if err != nil {
		return fmt.Errorf("get due recurring tasks: %w", err)
	}

	for _, parentTask := range tasks {
		if parentTask.RepeatUntil != nil && parentTask.RepeatUntil.Before(now) {
			parentTask.RepeatType = ""
			parentTask.RepeatConfig = nil
			parentTask.RepeatUntil = nil
			parentTask.NextOccurrence = nil
			_, _ = s.repo.Update(ctx, &parentTask)
			continue
		}

		newTask := &taskdomain.Task{
			Title:        parentTask.Title,
			Description:  parentTask.Description,
			Status:       taskdomain.StatusNew,
			CreatedAt:    now,
			UpdatedAt:    now,
			RepeatType:   parentTask.RepeatType,
			RepeatConfig: parentTask.RepeatConfig,
			ParentID:     &parentTask.ID,
		}

		_, err := s.repo.Create(ctx, newTask)
		if err != nil {
			fmt.Printf("failed to create recurring task for parent %d: %v\n", parentTask.ID, err)
			continue
		}

		nextDate := calculateNextOccurrence(now, &parentTask)

		parentTask.NextOccurrence = nextDate
		parentTask.UpdatedAt = now

		_, err = s.repo.Update(ctx, &parentTask)
		if err != nil {
			fmt.Printf("failed to update parent task %d next_occurrence: %v\n", parentTask.ID, err)
		}
	}

	return nil
}

func calculateNextOccurrence(now time.Time, task *taskdomain.Task) *time.Time {
	if task.RepeatType == "" {
		return nil
	}

	targetHour := now.Hour()
	targetMinute := now.Minute()

	if task.RepeatTime != "" {
		fmt.Sscanf(task.RepeatTime, "%d:%d", &targetHour, &targetMinute)
	}

	switch task.RepeatType {
	case "daily":
		var config struct {
			Interval int `json:"interval"`
		}
		_ = json.Unmarshal(task.RepeatConfig, &config)
		if config.Interval <= 0 {
			config.Interval = 1
		}
		next := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMinute, 0, 0, time.UTC)

		if next.Before(now) {
			next = next.AddDate(0, 0, config.Interval)
		} else if config.Interval > 1 {
			next = next.AddDate(0, 0, config.Interval-1)
		}
		return &next

	case "monthly":
		var config struct {
			DayOfMonth int `json:"day_of_month"`
		}
		_ = json.Unmarshal(task.RepeatConfig, &config)
		if config.DayOfMonth < 1 || config.DayOfMonth > 31 {
			config.DayOfMonth = 1
		}
		next := time.Date(now.Year(), now.Month()+1, config.DayOfMonth, targetHour, targetMinute, 0, 0, time.UTC)
		return &next

	case "parity":
		var config struct {
			Type string `json:"type"`
		}
		_ = json.Unmarshal(task.RepeatConfig, &config)

		next := now.AddDate(0, 0, 1)
		for {
			day := next.Day()
			if config.Type == "even" && day%2 == 0 {
				break
			}
			if config.Type == "odd" && day%2 == 1 {
				break
			}
			next = next.AddDate(0, 0, 1)
		}
		result := time.Date(next.Year(), next.Month(), next.Day(), targetHour, targetMinute, 0, 0, time.UTC)
		return &result

	case "specific_dates":
		var config struct {
			Dates []string `json:"dates"`
		}
		_ = json.Unmarshal(task.RepeatConfig, &config)

		today := now.Format("2006-01-02")
		for _, date := range config.Dates {
			if date > today {
				next, err := time.Parse("2006-01-02", date)
				if err != nil {
					return nil
				}
				result := time.Date(next.Year(), next.Month(), next.Day(), targetHour, targetMinute, 0, 0, time.UTC)
				return &result
			}
		}
		return nil
	}

	return nil
}
