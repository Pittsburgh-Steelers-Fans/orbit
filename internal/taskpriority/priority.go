package taskpriority

import (
	"context"
	"strings"
)

// Priority describes the urgency of a task.
type Priority string

const (
	// PriorityLow is for non-urgent tasks.
	PriorityLow Priority = "low"
	// PriorityMedium is for normal tasks.
	PriorityMedium Priority = "medium"
	// PriorityHigh is for urgent tasks.
	PriorityHigh Priority = "high"
)

// Task is a minimal task projection used for priority filtering.
type Task struct {
	ID       string   `json:"id"`
	Priority Priority `json:"priority"`
}

// Parse normalizes a priority string.
func Parse(ctx context.Context, value string) (Priority, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, err
	}
	priority := Priority(strings.ToLower(strings.TrimSpace(value)))
	return priority, priority.Valid(), nil
}

// Valid reports whether a priority is recognized.
func (p Priority) Valid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

// FilterByPriority returns tasks matching the requested priority.
func FilterByPriority(ctx context.Context, tasks []Task, priority Priority) ([]Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !priority.Valid() {
		return []Task{}, nil
	}
	filtered := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if task.Priority == priority {
			filtered = append(filtered, task)
		}
	}
	return filtered, nil
}
