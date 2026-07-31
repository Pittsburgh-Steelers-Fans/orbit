package projects

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrInvalidProject is returned when project input is invalid.
var ErrInvalidProject = errors.New("invalid project")

// Service coordinates project validation and persistence.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService creates a project service.
func NewService(repo Repository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{repo: repo, now: now}
}

// Create validates and creates a project.
func (s *Service) Create(ctx context.Context, name, ownerID string, members []string) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}
	name = strings.TrimSpace(name)
	ownerID = strings.TrimSpace(ownerID)
	if name == "" || ownerID == "" {
		return Project{}, ErrInvalidProject
	}
	createdAt := s.now().UTC()
	project := Project{ID: fmt.Sprintf("project-%d", createdAt.UnixNano()), Name: name, OwnerID: ownerID, Members: cleanMembers(ownerID, members), CreatedAt: createdAt}
	return s.repo.Create(ctx, project)
}

// Get returns a project by ID.
func (s *Service) Get(ctx context.Context, id string) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}
	if strings.TrimSpace(id) == "" {
		return Project{}, ErrInvalidProject
	}
	return s.repo.Get(ctx, id)
}

// Update changes a project's name and members.
func (s *Service) Update(ctx context.Context, id, name string, members []string) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}
	name = strings.TrimSpace(name)
	if strings.TrimSpace(id) == "" || name == "" {
		return Project{}, ErrInvalidProject
	}
	project, err := s.repo.Get(ctx, id)
	if err != nil {
		return Project{}, err
	}
	project.Name = name
	project.Members = cleanMembers(project.OwnerID, members)
	return s.repo.Update(ctx, project)
}

// Delete removes a project.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(id) == "" {
		return ErrInvalidProject
	}
	return s.repo.Delete(ctx, id)
}

// ListProjectsByUser returns projects visible to a user.
func (s *Service) ListProjectsByUser(ctx context.Context, userID string) ([]Project, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrInvalidProject
	}
	return s.repo.ListByUser(ctx, userID)
}

func cleanMembers(ownerID string, members []string) []string {
	seen := map[string]bool{ownerID: true}
	clean := make([]string, 0, len(members))
	for _, member := range members {
		member = strings.TrimSpace(member)
		if member == "" || seen[member] {
			continue
		}
		seen[member] = true
		clean = append(clean, member)
	}
	return clean
}
