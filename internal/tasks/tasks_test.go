package tasks

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceCreateDefaultsStatusAndSetsID(t *testing.T) {
	now := time.Date(2026, 7, 31, 18, 30, 0, 0, time.UTC)
	service := NewService(NewMemoryRepository(), func() time.Time { return now })

	task, err := service.Create(context.Background(), Task{ProjectID: "project-1", Title: " Ship API "})

	require.NoError(t, err)
	assert.Equal(t, "task-1785522600000000000", task.ID)
	assert.Equal(t, "Ship API", task.Title)
	assert.Equal(t, StatusTodo, task.Status)
	assert.Equal(t, now, task.CreatedAt)
}

func TestServiceCreateRejectsInvalidInput(t *testing.T) {
	service := NewService(NewMemoryRepository(), time.Now)

	cases := []Task{
		{ProjectID: "", Title: "Ship API"},
		{ProjectID: "project-1", Title: ""},
		{ProjectID: "project-1", Title: "Ship API", Status: Status("blocked")},
	}

	for _, tc := range cases {
		_, err := service.Create(context.Background(), tc)
		assert.ErrorIs(t, err, ErrInvalidTask)
	}
}

func TestServiceUpdateAndListTasks(t *testing.T) {
	service := NewService(NewMemoryRepository(), time.Now)
	task, err := service.Create(context.Background(), Task{ProjectID: "project-1", Title: "Ship API"})
	require.NoError(t, err)
	_, err = service.Create(context.Background(), Task{ProjectID: "project-2", Title: "Hide me"})
	require.NoError(t, err)
	dueDate := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

	updated, err := service.Update(context.Background(), task.ID, Task{Title: "Review API", Status: StatusReview, AssigneeID: " user-1 ", DueDate: &dueDate, Priority: "high"})
	require.NoError(t, err)
	assert.Equal(t, "Review API", updated.Title)
	assert.Equal(t, "user-1", updated.AssigneeID)

	tasks, err := service.List(context.Background(), "project-1")
	require.NoError(t, err)
	require.Len(t, tasks, 1)
	assert.Equal(t, task.ID, tasks[0].ID)
}

func TestServiceDeleteTask(t *testing.T) {
	service := NewService(NewMemoryRepository(), time.Now)
	task, err := service.Create(context.Background(), Task{ProjectID: "project-1", Title: "Ship API"})
	require.NoError(t, err)

	require.NoError(t, service.Delete(context.Background(), task.ID))
	_, err = service.Get(context.Background(), task.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRepositoryClonesDueDate(t *testing.T) {
	repo := NewMemoryRepository()
	dueDate := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	stored, err := repo.Create(context.Background(), Task{ID: "task-1", ProjectID: "project-1", Title: "Ship", DueDate: &dueDate})
	require.NoError(t, err)
	stored.DueDate = nil

	found, err := repo.Get(context.Background(), "task-1")
	require.NoError(t, err)
	require.NotNil(t, found.DueDate)
	assert.Equal(t, dueDate, *found.DueDate)
}
