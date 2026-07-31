package tasks

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler exposes task HTTP routes.
type Handler struct {
	service *Service
}

// NewHandler creates a task handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes returns the task router.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/tasks", h.createTask)
	r.Get("/tasks", h.listTasks)
	r.Get("/tasks/{taskID}", h.getTask)
	r.Put("/tasks/{taskID}", h.updateTask)
	r.Delete("/tasks/{taskID}", h.deleteTask)
	return r
}

// ServeHTTP serves task routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Routes().ServeHTTP(w, r)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	created, err := h.service.Create(r.Context(), task)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.List(r.Context(), r.URL.Query().Get("project_id"))
	if err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.service.Get(r.Context(), chi.URLParam(r, "taskID"))
	if err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	updated, err := h.service.Update(r.Context(), chi.URLParam(r, "taskID"), task)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), chi.URLParam(r, "taskID")); err != nil {
		writeTaskError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeTaskError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidTask):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
