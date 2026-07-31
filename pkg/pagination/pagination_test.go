package pagination

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCursorRoundTrip(t *testing.T) {
	for _, offset := range []int{0, 1, 25, 1024} {
		cursor := EncodeCursor(offset)
		require.NotEmpty(t, cursor)

		decoded, err := DecodeCursor(cursor)
		require.NoError(t, err)
		require.Equal(t, offset, decoded)
	}
}

func TestDecodeCursorRejectsInvalidInput(t *testing.T) {
	_, err := DecodeCursor("not-base64")
	require.Error(t, err)
}

func TestNewPage(t *testing.T) {
	page := NewPage([]string{"alpha", "beta"}, "next")

	require.Equal(t, []string{"alpha", "beta"}, page.Items)
	require.Equal(t, "next", page.NextCursor)
}
