package taskassignee

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveAssigneeReturnsNilForNilAssigneePointer(t *testing.T) {
	assignee, err := ResolveAssignee(context.Background(), Task{ID: "task-1", AssigneeID: nil})

	require.NoError(t, err)
	assert.Nil(t, assignee)
}

func TestResolveAssigneeReturnsNilForEmptyAssignee(t *testing.T) {
	empty := "   "

	assignee, err := ResolveAssignee(context.Background(), Task{ID: "task-1", AssigneeID: &empty})

	require.NoError(t, err)
	assert.Nil(t, assignee)
}

func TestResolveAssigneeReturnsTrimmedAssignee(t *testing.T) {
	id := " user-1 "

	assignee, err := ResolveAssignee(context.Background(), Task{ID: "task-1", AssigneeID: &id})

	require.NoError(t, err)
	require.NotNil(t, assignee)
	assert.Equal(t, "user-1", assignee.ID)
}

func TestResolveAssigneeHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assignee, err := ResolveAssignee(ctx, Task{ID: "task-1"})

	assert.Nil(t, assignee)
	assert.ErrorIs(t, err, context.Canceled)
}
