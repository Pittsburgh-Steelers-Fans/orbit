package taskassignee

import (
	"context"
	"strings"
)

// Task contains the assignee pointer used by assignee resolution.
type Task struct {
	ID         string  `json:"id"`
	AssigneeID *string `json:"assignee_id"`
}

// Assignee describes the user assigned to a task.
type Assignee struct {
	ID string `json:"id"`
}

// ResolveAssignee returns the assigned user, if any.
// It fixes the historical nil-pointer bug where unassigned tasks dereferenced AssigneeID.
func ResolveAssignee(ctx context.Context, task Task) (*Assignee, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if task.AssigneeID == nil {
		return nil, nil
	}
	id := strings.TrimSpace(*task.AssigneeID)
	if id == "" {
		return nil, nil
	}
	return &Assignee{ID: id}, nil
}
