package taskstatus

import (
	"context"
	"errors"
)

// Status describes a task workflow state.
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

// ErrInvalidTransition is returned when a status transition is not allowed.
var ErrInvalidTransition = errors.New("invalid status transition")

// AllowedTransitions lists permitted task status moves.
var AllowedTransitions = map[Status][]Status{
	StatusTodo:       {StatusInProgress},
	StatusInProgress: {StatusReview, StatusTodo},
	StatusReview:     {StatusDone, StatusInProgress},
	StatusDone:       {StatusReview},
}

// CanTransition reports whether a task can move between statuses.
func CanTransition(from, to Status) bool {
	for _, allowed := range AllowedTransitions[from] {
		if allowed == to {
			return true
		}
	}
	return false
}

// Transition validates and returns the next status.
func Transition(ctx context.Context, from, to Status) (Status, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if !CanTransition(from, to) {
		return "", ErrInvalidTransition
	}
	return to, nil
}
