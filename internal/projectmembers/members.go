package projectmembers

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
)

// Membership connects a user to a project with a role.
type Membership struct {
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
}

// ErrInvalidMembership is returned when membership input is invalid.
var ErrInvalidMembership = errors.New("invalid membership")

// ErrMembershipNotFound is returned when a membership cannot be found.
var ErrMembershipNotFound = errors.New("membership not found")

// Service manages project memberships in memory.
type Service struct {
	mu      sync.RWMutex
	members map[string]map[string]Membership
}

// NewService creates an empty membership service.
func NewService() *Service {
	return &Service{members: make(map[string]map[string]Membership)}
}

// AddMember adds or replaces a project membership.
func (s *Service) AddMember(ctx context.Context, projectID, userID, role string) (Membership, error) {
	if err := ctx.Err(); err != nil {
		return Membership{}, err
	}
	projectID = strings.TrimSpace(projectID)
	userID = strings.TrimSpace(userID)
	role = strings.TrimSpace(role)
	if projectID == "" || userID == "" || role == "" {
		return Membership{}, ErrInvalidMembership
	}
	membership := Membership{ProjectID: projectID, UserID: userID, Role: role}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.members[projectID] == nil {
		s.members[projectID] = make(map[string]Membership)
	}
	s.members[projectID][userID] = membership
	return membership, nil
}

// RemoveMember removes a user from a project.
func (s *Service) RemoveMember(ctx context.Context, projectID, userID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	projectID = strings.TrimSpace(projectID)
	userID = strings.TrimSpace(userID)
	if projectID == "" || userID == "" {
		return ErrInvalidMembership
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	projectMembers := s.members[projectID]
	if projectMembers == nil {
		return ErrMembershipNotFound
	}
	if _, ok := projectMembers[userID]; !ok {
		return ErrMembershipNotFound
	}
	delete(projectMembers, userID)
	return nil
}

// ListMembers returns all memberships for a project.
func (s *Service) ListMembers(ctx context.Context, projectID string) ([]Membership, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return nil, ErrInvalidMembership
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	members := make([]Membership, 0, len(s.members[projectID]))
	for _, membership := range s.members[projectID] {
		members = append(members, membership)
	}
	sort.Slice(members, func(i, j int) bool { return members[i].UserID < members[j].UserID })
	return members, nil
}
