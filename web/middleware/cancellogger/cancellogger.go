package cancellogger

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"net/http"
)

// ConstructContextCancellationLoggerMiddleware builds a middleware for logging context cancellations.
//
// This allows easier debugging of why and where in the stack a context was closed.
func ConstructContextCancellationLoggerMiddleware(description string) func(http.Handler) http.Handler {
	middleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context() // next might change it so must remember it here

			next.ServeHTTP(w, r)

			cause := context.Cause(ctx)
			if cause != nil {
				aulogging.Logger.NoCtx().Info().WithErr(cause).Printf("context '%s' is closed: %s", description, cause.Error())
			}
		}
		return http.HandlerFunc(fn)
	}
	return middleware
}
