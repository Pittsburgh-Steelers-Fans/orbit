package taskdue

import "time"

// Task is the minimal task shape needed for due-date calculations.
type Task struct {
	ID      string    `json:"id"`
	DueDate time.Time `json:"due_date"`
}

// IsOverdue reports whether due is strictly before now.
func IsOverdue(due, now time.Time) bool {
	return !due.IsZero() && due.Before(now)
}
