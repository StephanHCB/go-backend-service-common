package security

import (
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"net/http"
	"regexp"
)

type AuthRequiredMiddlewareOptions struct {
	// AllowUnauthorized is the explicit list of method + url path combinations that can allow
	// unauthorized access. Allows regular expressions.
	//
	// examples: "PUT /v1/info", "GET /swagger-ui.*"
	AllowUnauthorized []string
}

func AuthRequiredMiddleware(options AuthRequiredMiddlewareOptions) func(http.Handler) http.Handler {
	allowRegexes := make([]*regexp.Regexp, 0)
	for _, pattern := range options.AllowUnauthorized {
		fullMatchPattern := "^" + pattern + "$"
		re, err := regexp.Compile(fullMatchPattern)
		if err != nil {
			aulogging.Logger.NoCtx().Error().WithErr(err).Printf("allow unauthorized regexp pattern '%s' had errors - pattern skipped - please fix your configuration - continuing: %s", fullMatchPattern, err.Error())
		} else {
			allowRegexes = append(allowRegexes, re)
		}
	}

	mw := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			claims := GetClaims(ctx)
			if claims != nil {
				next.ServeHTTP(w, r)
				return
			}

			actualRequest := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			allowThrough := false
			for _, re := range allowRegexes {
				if re.MatchString(actualRequest) {
					allowThrough = true
				}
			}

			if allowThrough {
				next.ServeHTTP(w, r)
				return
			}

			unauthorizedErrorHandler(ctx, w, r, "Authorization required", Now())
		}
		return http.HandlerFunc(fn)
	}
	return mw
}
