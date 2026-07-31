package search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexSearchRanksSubstringAndTokenMatches(t *testing.T) {
	ctx := context.Background()
	index := NewIndex()
	require.NoError(t, index.Add(ctx, "task", "t1", "Write deployment checklist"))
	require.NoError(t, index.Add(ctx, "project", "p1", "Deployment platform migration"))
	require.NoError(t, index.Add(ctx, "task", "t2", "Fix profile avatar"))

	results, err := index.Search(ctx, "deploy check")
	require.NoError(t, err)

	require.Len(t, results, 2)
	assert.Equal(t, "task", results[0].Kind)
	assert.Equal(t, "t1", results[0].ID)
	assert.Greater(t, results[0].Score, results[1].Score)
	assert.Equal(t, "project", results[1].Kind)
}

func TestIndexAddReplacesExistingRecord(t *testing.T) {
	ctx := context.Background()
	index := NewIndex()
	require.NoError(t, index.Add(ctx, "task", "t1", "old title"))
	require.NoError(t, index.Add(ctx, "task", "t1", "new launch title"))

	results, err := index.Search(ctx, "launch")
	require.NoError(t, err)

	require.Len(t, results, 1)
	assert.Equal(t, "new launch title", results[0].Text)
}

func TestHandlerReturnsSearchResults(t *testing.T) {
	ctx := context.Background()
	index := NewIndex()
	require.NoError(t, index.Add(ctx, "project", "p1", "Roadmap planning"))
	handler := NewHandler(index).Routes()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/search?q=roadmap", nil)
	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	var results []Result
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&results))
	require.Len(t, results, 1)
	assert.Equal(t, "p1", results[0].ID)
}
