package projects

import (
	"context"
	"errors"
	"sort"
	"sync"
)

// ErrNotFound is returned when a project cannot be found.
var ErrNotFound = errors.New("project not found")

// Repository stores projects.
type Repository interface {
	Create(context.Context, Project) (Project, error)
	Get(context.Context, string) (Project, error)
	Update(context.Context, Project) (Project, error)
	Delete(context.Context, string) error
	ListByUser(context.Context, string) ([]Project, error)
}

// MemoryRepository stores projects in memory for tests and local development.
type MemoryRepository struct {
	mu       sync.RWMutex
	projects map[string]Project
}

// NewMemoryRepository creates an empty project repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{projects: make(map[string]Project)}
}

// Create stores a project.
func (r *MemoryRepository) Create(ctx context.Context, project Project) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.projects[project.ID] = cloneProject(project)
	return cloneProject(project), nil
}

// Get returns a project by ID.
func (r *MemoryRepository) Get(ctx context.Context, id string) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	project, ok := r.projects[id]
	if !ok {
		return Project{}, ErrNotFound
	}
	return cloneProject(project), nil
}

// Update replaces an existing project.
func (r *MemoryRepository) Update(ctx context.Context, project Project) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.projects[project.ID]; !ok {
		return Project{}, ErrNotFound
	}
	r.projects[project.ID] = cloneProject(project)
	return cloneProject(project), nil
}

// Delete removes a project by ID.
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.projects[id]; !ok {
		return ErrNotFound
	}
	delete(r.projects, id)
	return nil
}

// ListByUser returns projects owned by or shared with a user.
func (r *MemoryRepository) ListByUser(ctx context.Context, userID string) ([]Project, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	projects := make([]Project, 0)
	for _, project := range r.projects {
		if project.OwnerID == userID || containsMember(project.Members, userID) {
			projects = append(projects, cloneProject(project))
		}
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].CreatedAt.Before(projects[j].CreatedAt) })
	return projects, nil
}

func cloneProject(project Project) Project {
	project.Members = append([]string(nil), project.Members...)
	return project
}

func containsMember(members []string, userID string) bool {
	for _, member := range members {
		if member == userID {
			return true
		}
	}
	return false
}
