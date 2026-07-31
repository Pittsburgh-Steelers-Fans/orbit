package taskstatus

import (
	"encoding/json"
	"errors"
	"net/http"
)

type transitionRequest struct {
	From Status `json:"from"`
	To   Status `json:"to"`
}

type transitionResponse struct {
	Status Status `json:"status"`
}

// TransitionHandler validates POST requests that change a task status.
func TransitionHandler(w http.ResponseWriter, r *http.Request) {
	var req transitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	status, err := Transition(r.Context(), req.From, req.To)
	if err != nil {
		if errors.Is(err, ErrInvalidTransition) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, transitionResponse{Status: status})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
