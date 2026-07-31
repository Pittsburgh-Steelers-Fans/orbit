package authexpiry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateExpiry(t *testing.T) {
	now := time.Date(2026, 7, 31, 14, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		exp     time.Time
		wantErr bool
	}{
		{name: "expired token rejected", exp: now.Add(-time.Second), wantErr: true},
		{name: "valid token accepted", exp: now.Add(time.Minute)},
		{name: "exactly now rejected", exp: now, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExpiry(tt.exp, now)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.Is(err, ErrExpired))
				return
			}
			require.NoError(t, err)
		})
	}
}
