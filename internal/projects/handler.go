package projects

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler exposes project HTTP routes.
type Handler struct {
	service *Service
}

// NewHandler creates a project handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes returns the project router.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/projects", h.createProject)
	r.Get("/projects", h.listProjects)
	r.Get("/projects/{projectID}", h.getProject)
	r.Put("/projects/{projectID}", h.updateProject)
	r.Delete("/projects/{projectID}", h.deleteProject)
	return r
}

// ServeHTTP serves project routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Routes().ServeHTTP(w, r)
}

type projectRequest struct {
	Name    string   `json:"name"`
	OwnerID string   `json:"owner_id"`
	Members []string `json:"members"`
}

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var req projectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	project, err := h.service.Create(r.Context(), req.Name, req.OwnerID, req.Members)
	if err != nil {
		writeProjectError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, project)
}

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.ListProjectsByUser(r.Context(), r.URL.Query().Get("user_id"))
	if err != nil {
		writeProjectError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	project, err := h.service.Get(r.Context(), chi.URLParam(r, "projectID"))
	if err != nil {
		writeProjectError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	var req projectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	project, err := h.service.Update(r.Context(), chi.URLParam(r, "projectID"), req.Name, req.Members)
	if err != nil {
		writeProjectError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), chi.URLParam(r, "projectID")); err != nil {
		writeProjectError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeProjectError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidProject):
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
