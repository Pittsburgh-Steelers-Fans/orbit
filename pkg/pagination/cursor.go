package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

type cursorPayload struct {
	Offset int `json:"offset"`
}

// EncodeCursor returns an opaque cursor for the provided offset.
func EncodeCursor(offset int) string {
	payload, err := json.Marshal(cursorPayload{Offset: offset})
	if err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(payload)
}

// DecodeCursor parses an opaque cursor produced by EncodeCursor.
func DecodeCursor(s string) (int, error) {
	if s == "" {
		return 0, errors.New("cursor is empty")
	}

	payload, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, fmt.Errorf("decode cursor: %w", err)
	}

	var decoded cursorPayload
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return 0, fmt.Errorf("unmarshal cursor: %w", err)
	}
	if decoded.Offset < 0 {
		return 0, errors.New("cursor offset is negative")
	}

	return decoded.Offset, nil
}
