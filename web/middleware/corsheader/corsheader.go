package corsheader

import (
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
	"github.com/StephanHCB/go-backend-service-common/web/middleware/requestid"
	"github.com/go-http-utils/headers"
	"net/http"
	"strings"
)

// CorsHandling is deprecated and only retained for backwards compatibility.
//
// Please use CorsHandlingWithCorsAllowOrigin() or CorsHandlingWithConfig() instead.
func CorsHandling(next http.Handler) http.Handler {
	return CorsHandlingWithCorsAllowOrigin("*")(next)
}

// CorsHandlingWithConfig creates a middleware for CORS headers.
//
// It uses the configuration to decide whether to create an empty middleware or
// a middleware that actually sends disable CORS headers.
//
// With this it is now possible to control the precise contents of the
//
//	"Access-Control-Allow-Origin"
//
// header through this configuration values. If you need the service to forward
// authorization headers, you will need to limit yourself to a single base URL.
// If you do not need authorization forwarding, "*" is fine.
//
// Leaving the configuration value at its default of "" will switch off the
// CORS middleware completely.
func CorsHandlingWithConfig(configuration repository.Configuration) func(http.Handler) http.Handler {
	allowOrigin := configuration.CorsAllowOrigin()
	return CorsHandlingWithCorsAllowOrigin(allowOrigin)
}

// CorsHandlingWithCorsAllowOrigin creates a middleware for CORS headers.
//
// It uses the allowOrigin parameter to decide whether to create an empty middleware or
// a middleware that actually sends disable CORS headers.
//
// If you need the service to forward authorization headers, you will need to limit
// yourself to a single base URL.
//
// If you do not need authorization forwarding, "*" is fine.
//
// If allowOrigin is "", an empty middleware is returned.
func CorsHandlingWithCorsAllowOrigin(allowOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if allowOrigin == "" {
			// return an empty handler
			fn := func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		} else {
			// return a handler that sends appropriate cors disable headers
			fn := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(headers.AccessControlAllowOrigin, allowOrigin)

				w.Header().Set(headers.AccessControlAllowMethods, strings.Join([]string{
					http.MethodGet,
					http.MethodHead,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
				}, ", "))

				w.Header().Set(headers.AccessControlAllowHeaders, strings.Join([]string{
					headers.Accept,
					headers.ContentType,
					requestid.RequestIDHeader,
				}, ", "))

				w.Header().Set(headers.AccessControlAllowCredentials, "true")

				w.Header().Set(headers.AccessControlExposeHeaders, strings.Join([]string{
					headers.CacheControl,
					headers.ContentSecurityPolicy,
					headers.ContentType,
					headers.Location,
					requestid.RequestIDHeader,
				}, ", "))

				if r.Method == http.MethodOptions {
					// respond with ok, don't pass down the chain
					w.WriteHeader(http.StatusOK)
					return
				}

				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		}
	}
}
