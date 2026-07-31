package taskstatus

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanTransitionAllowsConfiguredTransitions(t *testing.T) {
	cases := []struct {
		from Status
		to   Status
	}{
		{StatusTodo, StatusInProgress},
		{StatusInProgress, StatusReview},
		{StatusReview, StatusDone},
		{StatusDone, StatusReview},
	}

	for _, tc := range cases {
		assert.True(t, CanTransition(tc.from, tc.to))
	}
}

func TestTransitionRejectsInvalidMoves(t *testing.T) {
	_, err := Transition(context.Background(), StatusTodo, StatusDone)

	assert.ErrorIs(t, err, ErrInvalidTransition)
}

func TestTransitionHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Transition(ctx, StatusTodo, StatusInProgress)

	assert.ErrorIs(t, err, context.Canceled)
}

func TestTransitionHandlerReturnsNextStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/tasks/task-1/status", bytes.NewBufferString(`{"from":"todo","to":"in_progress"}`))
	rec := httptest.NewRecorder()

	TransitionHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"in_progress"}`, rec.Body.String())
}

func TestTransitionHandlerRejectsInvalidTransition(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/tasks/task-1/status", bytes.NewBufferString(`{"from":"todo","to":"done"}`))
	rec := httptest.NewRecorder()

	TransitionHandler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"invalid status transition"}`, rec.Body.String())
}
