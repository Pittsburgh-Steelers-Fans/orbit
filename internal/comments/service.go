package comments

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ErrNotFound is returned when a comment cannot be found.
var ErrNotFound = errors.New("comment not found")

// Repository stores and retrieves comments.
type Repository interface {
	Create(ctx context.Context, comment Comment) (Comment, error)
	List(ctx context.Context, taskID string) ([]Comment, error)
	Delete(ctx context.Context, id string) error
}

// MemoryRepository is a concurrency-safe in-memory comment repository.
type MemoryRepository struct {
	mu       sync.RWMutex
	comments map[string]Comment
	nextID   int
}

// NewMemoryRepository creates an empty in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{comments: make(map[string]Comment)}
}

// Create stores comment and assigns an ID when one is not provided.
func (r *MemoryRepository) Create(ctx context.Context, comment Comment) (Comment, error) {
	if err := ctx.Err(); err != nil {
		return Comment{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if comment.ID == "" {
		r.nextID++
		comment.ID = fmt.Sprintf("comment-%d", r.nextID)
	}
	r.comments[comment.ID] = comment
	return comment, nil
}

// List returns comments for taskID in insertion order by generated ID.
func (r *MemoryRepository) List(ctx context.Context, taskID string) ([]Comment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	results := make([]Comment, 0)
	for i := 1; i <= r.nextID; i++ {
		comment, ok := r.comments[fmt.Sprintf("comment-%d", i)]
		if ok && comment.TaskID == taskID {
			results = append(results, comment)
		}
	}
	return results, nil
}

// Delete removes a comment by ID.
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.comments[id]; !ok {
		return ErrNotFound
	}
	delete(r.comments, id)
	return nil
}

// Service coordinates comment validation and persistence.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService creates a comment service backed by repo.
func NewService(repo Repository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{repo: repo, now: now}
}

// Create validates and stores a comment.
func (s *Service) Create(ctx context.Context, taskID, authorID, body string) (Comment, error) {
	comment := Comment{
		TaskID:    strings.TrimSpace(taskID),
		AuthorID:  strings.TrimSpace(authorID),
		Body:      strings.TrimSpace(body),
		CreatedAt: s.now().UTC(),
	}
	if comment.TaskID == "" || comment.AuthorID == "" || comment.Body == "" {
		return Comment{}, errors.New("task_id, author_id, and body are required")
	}
	return s.repo.Create(ctx, comment)
}

// List returns comments for a task.
func (s *Service) List(ctx context.Context, taskID string) ([]Comment, error) {
	return s.repo.List(ctx, strings.TrimSpace(taskID))
}

// Delete removes a comment by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, strings.TrimSpace(id))
}
