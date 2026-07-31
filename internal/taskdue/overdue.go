package taskdue

import (
	"encoding/json"
	"net/http"
	"time"
)

// ListOverdue returns tasks whose due date is before now, preserving input order.
func ListOverdue(tasks []Task, now time.Time) []Task {
	overdue := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if IsOverdue(task.DueDate, now) {
			overdue = append(overdue, task)
		}
	}
	return overdue
}

// OverdueHandler returns an HTTP handler that writes the current overdue task list as JSON.
func OverdueHandler(tasks []Task, now func() time.Time) http.HandlerFunc {
	if now == nil {
		now = time.Now
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Tasks []Task `json:"tasks"`
		}{Tasks: ListOverdue(tasks, now())})
	}
}
