package health

import (
	"context"
	"net/http"
)

// Check verifies whether one readiness dependency is available.
type Check func(ctx context.Context) error

// AddCheck registers another readiness dependency check.
func (h *Handler) AddCheck(check Check) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks = append(h.checks, check)
}

// Ready returns ok when every registered dependency check succeeds.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if err := r.Context().Err(); err != nil {
		writeStatus(w, http.StatusRequestTimeout, map[string]string{"status": "error", "error": err.Error()})
		return
	}
	h.mu.RLock()
	checks := append([]Check(nil), h.checks...)
	h.mu.RUnlock()

	for _, check := range checks {
		if err := check(r.Context()); err != nil {
			writeStatus(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable", "error": err.Error()})
			return
		}
	}
	writeStatus(w, http.StatusOK, map[string]string{"status": "ok"})
}
