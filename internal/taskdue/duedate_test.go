package taskdue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIsOverdue(t *testing.T) {
	now := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)

	require.True(t, IsOverdue(now.Add(-time.Minute), now))
	require.False(t, IsOverdue(now, now))
	require.False(t, IsOverdue(now.Add(time.Minute), now))
	require.False(t, IsOverdue(time.Time{}, now))
}

func TestListOverdue(t *testing.T) {
	now := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	tasks := []Task{
		{ID: "future", DueDate: now.Add(time.Hour)},
		{ID: "past", DueDate: now.Add(-time.Hour)},
		{ID: "older", DueDate: now.Add(-2 * time.Hour)},
	}

	overdue := ListOverdue(tasks, now)

	require.Equal(t, []Task{tasks[1], tasks[2]}, overdue)
}

func TestOverdueHandler(t *testing.T) {
	now := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	handler := OverdueHandler([]Task{
		{ID: "past", DueDate: now.Add(-time.Hour)},
		{ID: "future", DueDate: now.Add(time.Hour)},
	}, func() time.Time { return now })

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/overdue", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	var body struct {
		Tasks []Task `json:"tasks"`
	}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Len(t, body.Tasks, 1)
	require.Equal(t, "past", body.Tasks[0].ID)
}
