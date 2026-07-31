package paginationfix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeExclusiveDoesNotReturnNextPageSentinel(t *testing.T) {
	rows := []string{"r0", "r1", "r2", "r3", "r4", "r5", "r6", "r7", "r8", "r9", "r10"}

	start, end, err := DecodeExclusive("4", 5)
	require.NoError(t, err)
	require.Equal(t, 5, start)
	require.Equal(t, 10, end)
	require.Equal(t, []string{"r5", "r6", "r7", "r8", "r9"}, rows[start:end])
	require.NotContains(t, rows[start:end], "r10")
}

func TestDecodeExclusiveFirstPage(t *testing.T) {
	start, end, err := DecodeExclusive("", 3)
	require.NoError(t, err)
	require.Equal(t, 0, start)
	require.Equal(t, 3, end)
}
