package users

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Service coordinates user validation and persistence.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService creates a user service.
func NewService(repo Repository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{repo: repo, now: now}
}

// ListUsers returns all users.
func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}

// GetUser returns one user by id.
func (s *Service) GetUser(ctx context.Context, id string) (User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return User{}, errors.New("users: id is required")
	}
	return s.repo.Get(ctx, id)
}

// UpdateUser validates and updates an existing user.
func (s *Service) UpdateUser(ctx context.Context, id string, email string, name string) (User, error) {
	id = strings.TrimSpace(id)
	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)
	if id == "" || email == "" || name == "" {
		return User{}, errors.New("users: id, email, and name are required")
	}
	current, err := s.repo.Get(ctx, id)
	if err != nil {
		return User{}, err
	}
	current.Email = email
	current.Name = name
	return s.repo.Update(ctx, current)
}

// DeleteUser removes a user by id.
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("users: id is required")
	}
	return s.repo.Delete(ctx, id)
}
