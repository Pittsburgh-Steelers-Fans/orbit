package comments

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServiceCreateListDelete(t *testing.T) {
	now := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	service := NewService(NewMemoryRepository(), func() time.Time { return now })

	created, err := service.Create(t.Context(), "task-1", "user-1", " looks good ")
	require.NoError(t, err)
	require.Equal(t, "comment-1", created.ID)
	require.Equal(t, "looks good", created.Body)
	require.Equal(t, now, created.CreatedAt)

	comments, err := service.List(t.Context(), "task-1")
	require.NoError(t, err)
	require.Equal(t, []Comment{created}, comments)

	require.NoError(t, service.Delete(t.Context(), created.ID))
	comments, err = service.List(t.Context(), "task-1")
	require.NoError(t, err)
	require.Empty(t, comments)
}

func TestHandlerCreateListAndDelete(t *testing.T) {
	service := NewService(NewMemoryRepository(), func() time.Time {
		return time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	})
	handler := NewHandler(service)

	payload := bytes.NewBufferString(`{"task_id":"task-1","author_id":"user-1","body":"ship it"}`)
	createRecorder := httptest.NewRecorder()
	handler.Create(createRecorder, httptest.NewRequest(http.MethodPost, "/comments", payload))
	require.Equal(t, http.StatusCreated, createRecorder.Code)

	var created Comment
	require.NoError(t, json.NewDecoder(createRecorder.Body).Decode(&created))
	require.Equal(t, "comment-1", created.ID)

	listRecorder := httptest.NewRecorder()
	handler.List(listRecorder, httptest.NewRequest(http.MethodGet, "/comments?task_id=task-1", nil))
	require.Equal(t, http.StatusOK, listRecorder.Code)

	var listed []Comment
	require.NoError(t, json.NewDecoder(listRecorder.Body).Decode(&listed))
	require.Len(t, listed, 1)

	deleteRecorder := httptest.NewRecorder()
	handler.Delete(deleteRecorder, httptest.NewRequest(http.MethodDelete, "/comments/"+created.ID, nil))
	require.Equal(t, http.StatusNoContent, deleteRecorder.Code)
}
