package pagination

// CursorPage contains one page of items and the cursor for the next page.
type CursorPage[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// NewPage returns a cursor page with the provided items and next cursor.
func NewPage[T any](items []T, nextCursor string) CursorPage[T] {
	return CursorPage[T]{
		Items:      items,
		NextCursor: nextCursor,
	}
}
