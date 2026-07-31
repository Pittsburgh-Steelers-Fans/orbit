package users

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceUpdateUser(t *testing.T) {
	created := time.Date(2026, 7, 31, 14, 0, 0, 0, time.UTC)
	repo := NewMemoryRepository(User{ID: "user-1", Email: "old@example.com", Name: "Old", CreatedAt: created})
	service := NewService(repo, nil)

	user, err := service.UpdateUser(context.Background(), "user-1", "NEW@EXAMPLE.COM", "New Name")
	require.NoError(t, err)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, "New Name", user.Name)
	assert.Equal(t, created, user.CreatedAt)
}

func TestHandlerListAndDeleteUsers(t *testing.T) {
	repo := NewMemoryRepository(User{ID: "user-1", Email: "one@example.com", Name: "One", CreatedAt: time.Now()})
	handler := NewHandler(NewService(repo, nil))
	router := chi.NewRouter()
	handler.Routes(router)

	listReq := httptest.NewRequest(http.MethodGet, "/users", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)
	require.Equal(t, http.StatusOK, listRec.Code)
	var users []User
	require.NoError(t, json.NewDecoder(listRec.Body).Decode(&users))
	assert.Len(t, users, 1)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/users/user-1", nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteRec.Code)
}

func TestHandlerUpdateUser(t *testing.T) {
	repo := NewMemoryRepository(User{ID: "user-1", Email: "one@example.com", Name: "One", CreatedAt: time.Now()})
	handler := NewHandler(NewService(repo, nil))
	router := chi.NewRouter()
	handler.Routes(router)

	body := bytes.NewBufferString(`{"email":"two@example.com","name":"Two"}`)
	req := httptest.NewRequest(http.MethodPut, "/users/user-1", body)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var user User
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&user))
	assert.Equal(t, "two@example.com", user.Email)
	assert.Equal(t, "Two", user.Name)
}
