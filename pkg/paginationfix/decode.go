package paginationfix

import (
	"errors"
	"fmt"
	"strconv"
)

// DecodeExclusive converts a cursor into a half-open [start, end) range.
//
// The cursor represents the last row returned by the previous page. Older code
// treated the end index as inclusive while fetching pageSize+1 rows to detect a
// next page, which leaked the first row of the next page into the current page.
// Returning a half-open range keeps the extra N+1 sentinel row out of results.
func DecodeExclusive(cursor string, pageSize int) (start, end int, err error) {
	if pageSize < 1 {
		return 0, 0, errors.New("page size must be positive")
	}

	if cursor == "" {
		return 0, pageSize, nil
	}

	lastSeen, err := strconv.Atoi(cursor)
	if err != nil {
		return 0, 0, fmt.Errorf("decode cursor: %w", err)
	}
	if lastSeen < 0 {
		return 0, 0, errors.New("cursor must not be negative")
	}

	start = lastSeen + 1
	end = start + pageSize
	return start, end, nil
}
