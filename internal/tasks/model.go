package tasks

import "time"

// Status describes where a task is in the workflow.
type Status string

const (
	// StatusTodo means work has not started.
	StatusTodo Status = "todo"
	// StatusInProgress means work is underway.
	StatusInProgress Status = "in_progress"
	// StatusReview means work is waiting for review.
	StatusReview Status = "review"
	// StatusDone means work is complete.
	StatusDone Status = "done"
)

// Task captures a unit of work in a project.
type Task struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	AssigneeID  string     `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	Priority    string     `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
}
