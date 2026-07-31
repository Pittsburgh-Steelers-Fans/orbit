package search

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler serves search requests against an Index.
type Handler struct {
	index *Index
}

// NewHandler creates a search handler using the supplied index.
func NewHandler(index *Index) *Handler {
	return &Handler{index: index}
}

// Routes wires the search endpoint onto a chi router.
func (h *Handler) Routes() http.Handler {
	router := chi.NewRouter()
	router.Get("/search", h.Search)
	return router
}

// Search writes ranked search results for the q query parameter.
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	results, err := h.index.Search(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(results)
}
