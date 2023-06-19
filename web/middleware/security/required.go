package security

import (
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"net/http"
	"path/filepath"
)

type AuthRequiredMiddlewareOptions struct {
	// AllowUnauthorized is the explicit list of method + url path combinations that can allow
	// unauthorized access.
	//
	// examples: "PUT /v1/info", "GET /swagger-ui/*" (* glob supported in path)
	AllowUnauthorized []string
}

func AuthRequiredMiddleware(options AuthRequiredMiddlewareOptions) func(http.Handler) http.Handler {
	mw := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			actualRequest := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

			allowThrough := false
			for _, globPattern := range options.AllowUnauthorized {
				matched, err := filepath.Match(globPattern, actualRequest)
				if err != nil {
					aulogging.Logger.Ctx(ctx).Error().WithErr(err).Printf("allow unauthorized glob pattern '%s' had errors - pattern skipped - please fix your configuration: %s", globPattern, err.Error())
				} else {
					if matched {
						allowThrough = true
					}
				}
			}

			claims := GetClaims(ctx)
			if claims == nil && !allowThrough {
				unauthorizedErrorHandler(ctx, w, r, "Authorization required", Now())
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return mw
}
