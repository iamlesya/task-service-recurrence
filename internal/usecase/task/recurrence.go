package task

import (
	"encoding/json"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func CalculateNextOccurrence(now time.Time, task *taskdomain.Task) *time.Time {
	if task.RepeatType == "" {
		return nil
	}

	switch task.RepeatType {
	case "daily":
		var config struct {
			Interval int `json:"interval"`
		}
		if err := json.Unmarshal(task.RepeatConfig, &config); err != nil {
			return nil
		}
		if config.Interval <= 0 {
			config.Interval = 1
		}
		next := now.AddDate(0, 0, config.Interval)
		return &next

	case "monthly":
		var config struct {
			DayOfMonth int `json:"day_of_month"`
		}
		if err := json.Unmarshal(task.RepeatConfig, &config); err != nil {
			return nil
		}
		next := time.Date(now.Year(), now.Month()+1, config.DayOfMonth, 0, 0, 0, 0, time.UTC)
		return &next

		// TODO: добавить specific_dates и parity
	}

	return nil
}
