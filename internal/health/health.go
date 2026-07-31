package health

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Handler serves health and readiness probes.
type Handler struct {
	mu     sync.RWMutex
	checks []Check
}

// NewHandler creates a health handler with optional readiness checks.
func NewHandler(checks ...Check) *Handler {
	return &Handler{checks: append([]Check(nil), checks...)}
}

// Health returns a simple ok response for liveness probes.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if err := r.Context().Err(); err != nil {
		writeStatus(w, http.StatusRequestTimeout, map[string]string{"status": "error", "error": err.Error()})
		return
	}
	writeStatus(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeStatus(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
