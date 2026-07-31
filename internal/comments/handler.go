package comments

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Handler exposes HTTP CRUD endpoints for comments.
type Handler struct {
	service *Service
}

// NewHandler creates a Handler backed by service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST requests that create comments.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		TaskID   string `json:"task_id"`
		AuthorID string `json:"author_id"`
		Body     string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	comment, err := h.service.Create(r.Context(), input.TaskID, input.AuthorID, input.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

// List handles GET requests for comments belonging to a task.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	taskID := r.URL.Query().Get("task_id")
	comments, err := h.service.List(r.Context(), taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

// Delete handles DELETE requests and returns 204 No Content on success.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := commentIDFromRequest(r)
	if id == "" {
		http.Error(w, "comment id is required", http.StatusBadRequest)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func commentIDFromRequest(r *http.Request) string {
	if id := strings.TrimSpace(r.URL.Query().Get("id")); id != "" {
		return id
	}
	return strings.Trim(strings.TrimPrefix(r.URL.Path, "/comments/"), "/")
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
