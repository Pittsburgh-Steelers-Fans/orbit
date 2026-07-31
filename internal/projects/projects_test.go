package projects

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceCreateValidatesAndCleansMembers(t *testing.T) {
	now := time.Date(2026, 7, 31, 18, 0, 0, 0, time.UTC)
	service := NewService(NewMemoryRepository(), func() time.Time { return now })

	project, err := service.Create(context.Background(), " Launch ", "owner-1", []string{"owner-1", "user-2", "", "user-2"})

	require.NoError(t, err)
	assert.Equal(t, "project-1785520800000000000", project.ID)
	assert.Equal(t, "Launch", project.Name)
	assert.Equal(t, []string{"user-2"}, project.Members)
	assert.Equal(t, now, project.CreatedAt)
}

func TestServiceCreateRejectsMissingFields(t *testing.T) {
	service := NewService(NewMemoryRepository(), time.Now)

	cases := []struct {
		name    string
		ownerID string
	}{
		{name: "", ownerID: "owner-1"},
		{name: "Launch", ownerID: ""},
	}

	for _, tc := range cases {
		_, err := service.Create(context.Background(), tc.name, tc.ownerID, nil)
		assert.ErrorIs(t, err, ErrInvalidProject)
	}
}

func TestServiceListProjectsByUserIncludesOwnedAndMemberProjects(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo, time.Now)
	owned, err := service.Create(context.Background(), "Owned", "user-1", nil)
	require.NoError(t, err)
	shared, err := service.Create(context.Background(), "Shared", "user-2", []string{"user-1"})
	require.NoError(t, err)
	_, err = service.Create(context.Background(), "Hidden", "user-3", nil)
	require.NoError(t, err)

	projects, err := service.ListProjectsByUser(context.Background(), "user-1")

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{owned.ID, shared.ID}, []string{projects[0].ID, projects[1].ID})
}

func TestServiceUpdateAndDeleteProject(t *testing.T) {
	service := NewService(NewMemoryRepository(), time.Now)
	project, err := service.Create(context.Background(), "Launch", "owner-1", nil)
	require.NoError(t, err)

	updated, err := service.Update(context.Background(), project.ID, "Renamed", []string{"user-2"})
	require.NoError(t, err)
	assert.Equal(t, "Renamed", updated.Name)
	assert.Equal(t, []string{"user-2"}, updated.Members)

	require.NoError(t, service.Delete(context.Background(), project.ID))
	_, err = service.Get(context.Background(), project.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRepositoryReturnsContextError(t *testing.T) {
	repo := NewMemoryRepository()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.Create(ctx, Project{})

	assert.True(t, errors.Is(err, context.Canceled))
}
