package search

import (
	"context"
	"sort"
	"strings"
	"sync"
)

// Result represents one indexed item matched by a query.
type Result struct {
	Kind  string `json:"kind"`
	ID    string `json:"id"`
	Text  string `json:"text"`
	Score int    `json:"score"`
}

// Index is a concurrency-safe in-memory search index for Orbit resources.
type Index struct {
	mu      sync.RWMutex
	records []record
}

type record struct {
	kind string
	id   string
	text string
}

// NewIndex creates an empty in-memory index.
func NewIndex() *Index {
	return &Index{}
}

// Add inserts or replaces a searchable item.
func (i *Index) Add(ctx context.Context, kind, id, text string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	for pos, existing := range i.records {
		if existing.kind == kind && existing.id == id {
			i.records[pos] = record{kind: kind, id: id, text: text}
			return nil
		}
	}
	i.records = append(i.records, record{kind: kind, id: id, text: text})
	return nil
}

// Search returns ranked results using case-insensitive substring and token matches.
func (i *Index) Search(ctx context.Context, query string) ([]Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return []Result{}, nil
	}
	tokens := strings.Fields(query)

	i.mu.RLock()
	defer i.mu.RUnlock()
	results := make([]Result, 0, len(i.records))
	for _, entry := range i.records {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		score := scoreRecord(entry.text, query, tokens)
		if score > 0 {
			results = append(results, Result{Kind: entry.kind, ID: entry.id, Text: entry.text, Score: score})
		}
	}
	sort.Slice(results, func(left, right int) bool {
		if results[left].Score != results[right].Score {
			return results[left].Score > results[right].Score
		}
		if results[left].Kind != results[right].Kind {
			return results[left].Kind < results[right].Kind
		}
		return results[left].ID < results[right].ID
	})
	return results, nil
}

func scoreRecord(text, query string, tokens []string) int {
	lower := strings.ToLower(text)
	score := 0
	if strings.Contains(lower, query) {
		score += 10 + len(query)
	}
	for _, token := range tokens {
		if strings.Contains(lower, token) {
			score += 3 + len(token)
		}
	}
	return score
}
