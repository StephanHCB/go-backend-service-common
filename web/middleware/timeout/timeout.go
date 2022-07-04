package timeout

import (
	"context"
	"net/http"
	"time"
)

var RequestTimeoutSeconds = 30

func AddRequestTimeout(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx, cancel := context.WithTimeout(ctx, time.Duration(RequestTimeoutSeconds)*time.Second)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
