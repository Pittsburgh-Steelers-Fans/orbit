package tasks

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrInvalidTask is returned when task input is invalid.
var ErrInvalidTask = errors.New("invalid task")

// Service coordinates task validation and persistence.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService creates a task service.
func NewService(repo Repository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{repo: repo, now: now}
}

// Create validates and creates a task.
func (s *Service) Create(ctx context.Context, task Task) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	task.ProjectID = strings.TrimSpace(task.ProjectID)
	task.Title = strings.TrimSpace(task.Title)
	if task.ProjectID == "" || task.Title == "" {
		return Task{}, ErrInvalidTask
	}
	if task.Status == "" {
		task.Status = StatusTodo
	}
	if !validStatus(task.Status) {
		return Task{}, ErrInvalidTask
	}
	createdAt := s.now().UTC()
	task.ID = fmt.Sprintf("task-%d", createdAt.UnixNano())
	task.CreatedAt = createdAt
	return s.repo.Create(ctx, task)
}

// Get returns a task by ID.
func (s *Service) Get(ctx context.Context, id string) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	if strings.TrimSpace(id) == "" {
		return Task{}, ErrInvalidTask
	}
	return s.repo.Get(ctx, id)
}

// Update replaces mutable task fields.
func (s *Service) Update(ctx context.Context, id string, next Task) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	if strings.TrimSpace(id) == "" || strings.TrimSpace(next.Title) == "" {
		return Task{}, ErrInvalidTask
	}
	if next.Status == "" {
		next.Status = StatusTodo
	}
	if !validStatus(next.Status) {
		return Task{}, ErrInvalidTask
	}
	current, err := s.repo.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}
	current.Title = strings.TrimSpace(next.Title)
	current.Description = next.Description
	current.Status = next.Status
	current.AssigneeID = strings.TrimSpace(next.AssigneeID)
	current.DueDate = next.DueDate
	current.Priority = strings.TrimSpace(next.Priority)
	return s.repo.Update(ctx, current)
}

// Delete removes a task.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(id) == "" {
		return ErrInvalidTask
	}
	return s.repo.Delete(ctx, id)
}

// List returns tasks for a project.
func (s *Service) List(ctx context.Context, projectID string) ([]Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, strings.TrimSpace(projectID))
}

func validStatus(status Status) bool {
	switch status {
	case StatusTodo, StatusInProgress, StatusReview, StatusDone:
		return true
	default:
		return false
	}
}
