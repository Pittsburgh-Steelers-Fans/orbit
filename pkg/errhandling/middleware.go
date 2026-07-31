package errhandling

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
)

// Recover converts panics into JSON 500 responses and records the panic with the logger.
func Recover(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error().Bytes("stack", debug.Stack()).Interface("panic", recovered).Msg("panic recovered")
					Respond(w, Wrap(http.StatusInternalServerError, defaultErrorCode, "internal server error", fmt.Errorf("panic: %v", recovered)))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
