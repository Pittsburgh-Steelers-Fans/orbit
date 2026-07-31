package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler exposes user CRUD endpoints.
type Handler struct {
	service *Service
}

type updateRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// NewHandler creates a user handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes mounts user routes on a chi router.
func (h *Handler) Routes(r chi.Router) {
	r.Get("/users", h.List)
	r.Get("/users/{userID}", h.Get)
	r.Put("/users/{userID}", h.Update)
	r.Delete("/users/{userID}", h.Delete)
}

// List writes all users as JSON.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// Get writes a single user as JSON.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetUser(r.Context(), chi.URLParam(r, "userID"))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// Update replaces mutable user fields.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.service.UpdateUser(r.Context(), chi.URLParam(r, "userID"), req.Email, req.Name)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// Delete removes a user.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteUser(r.Context(), chi.URLParam(r, "userID")); err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		h.writeError(w, http.StatusNotFound, err)
		return
	}
	h.writeError(w, http.StatusBadRequest, err)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
