package taskpriority

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNormalizesKnownPriorities(t *testing.T) {
	priority, ok, err := Parse(context.Background(), " HIGH ")

	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, PriorityHigh, priority)
}

func TestParseRejectsUnknownPriority(t *testing.T) {
	priority, ok, err := Parse(context.Background(), "urgent")

	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, Priority("urgent"), priority)
}

func TestPriorityValidRecognizesAllowedValues(t *testing.T) {
	assert.True(t, PriorityLow.Valid())
	assert.True(t, PriorityMedium.Valid())
	assert.True(t, PriorityHigh.Valid())
	assert.False(t, Priority("urgent").Valid())
}

func TestFilterByPriorityReturnsMatchingTasks(t *testing.T) {
	tasks := []Task{
		{ID: "task-1", Priority: PriorityHigh},
		{ID: "task-2", Priority: PriorityLow},
		{ID: "task-3", Priority: PriorityHigh},
	}

	filtered, err := FilterByPriority(context.Background(), tasks, PriorityHigh)

	require.NoError(t, err)
	assert.Equal(t, []Task{{ID: "task-1", Priority: PriorityHigh}, {ID: "task-3", Priority: PriorityHigh}}, filtered)
}

func TestFilterByPriorityHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	filtered, err := FilterByPriority(ctx, nil, PriorityHigh)

	assert.Nil(t, filtered)
	assert.ErrorIs(t, err, context.Canceled)
}
