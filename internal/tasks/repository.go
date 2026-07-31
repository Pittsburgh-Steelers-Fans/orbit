package tasks

import (
	"context"
	"errors"
	"sort"
	"sync"
)

// ErrNotFound is returned when a task cannot be found.
var ErrNotFound = errors.New("task not found")

// Repository stores tasks.
type Repository interface {
	Create(context.Context, Task) (Task, error)
	Get(context.Context, string) (Task, error)
	Update(context.Context, Task) (Task, error)
	Delete(context.Context, string) error
	List(context.Context, string) ([]Task, error)
}

// MemoryRepository stores tasks in memory.
type MemoryRepository struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

// NewMemoryRepository creates an empty task repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{tasks: make(map[string]Task)}
}

// Create stores a task.
func (r *MemoryRepository) Create(ctx context.Context, task Task) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = cloneTask(task)
	return cloneTask(task), nil
}

// Get returns a task by ID.
func (r *MemoryRepository) Get(ctx context.Context, id string) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return cloneTask(task), nil
}

// Update replaces an existing task.
func (r *MemoryRepository) Update(ctx context.Context, task Task) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[task.ID]; !ok {
		return Task{}, ErrNotFound
	}
	r.tasks[task.ID] = cloneTask(task)
	return cloneTask(task), nil
}

// Delete removes a task.
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(r.tasks, id)
	return nil
}

// List returns tasks for a project, or all tasks if projectID is empty.
func (r *MemoryRepository) List(ctx context.Context, projectID string) ([]Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	tasks := make([]Task, 0)
	for _, task := range r.tasks {
		if projectID == "" || task.ProjectID == projectID {
			tasks = append(tasks, cloneTask(task))
		}
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].CreatedAt.Before(tasks[j].CreatedAt) })
	return tasks, nil
}

func cloneTask(task Task) Task {
	if task.DueDate != nil {
		dueDate := *task.DueDate
		task.DueDate = &dueDate
	}
	return task
}
