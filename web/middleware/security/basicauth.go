package security

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/go-http-utils/headers"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

type BasicAuthMiddlewareOptions struct {
	// supports a fixed basic auth setup for use by e.g. CI systems
	//
	// you must provide nonempty username and password from configuration,
	// when the authorization header starts with "Basic" and matches, the injected
	// user will then have the CustomClaims provided here
	BasicAuthUsername string
	BasicAuthPassword string
	BasicAuthClaims   CustomClaims
}

func BasicAuthValidatorMiddleware(options BasicAuthMiddlewareOptions) func(http.Handler) http.Handler {
	mw := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeaderValue := r.Header.Get(headers.Authorization)
			const basicPrefix = "Basic "
			if !strings.HasPrefix(authHeaderValue, basicPrefix) {
				// valid case, no basic auth authorization provided, fall through
				next.ServeHTTP(w, r)
			} else {
				ctx := r.Context()
				username, password, basicAuthOk := r.BasicAuth()
				if basicAuthOk {
					if checkBasicAuthValue(username, password, options) {
						adminClaims := AllClaims{
							RegisteredClaims: jwt.RegisteredClaims{},
							CustomClaims:     options.BasicAuthClaims,
						}
						ctx = PutClaims(ctx, &adminClaims)
						next.ServeHTTP(w, r.WithContext(ctx))
					} else {
						unauthorizedErrorHandler(ctx, w, r, "Authorization failed Basic Auth", Now())
						return
					}
				} else {
					unauthorizedErrorHandler(ctx, w, r, "Authorization header contains invalid values for basic auth", Now())
					return
				}
			}
		}
		return http.HandlerFunc(fn)
	}
	return mw
}

func checkBasicAuthValue(username string, password string, options BasicAuthMiddlewareOptions) bool {
	if username == "" || password == "" {
		return false
	}

	expectedUsernameHash := sha256.Sum256([]byte(options.BasicAuthUsername))
	expectedPasswordHash := sha256.Sum256([]byte(options.BasicAuthPassword))

	usernameHash := sha256.Sum256([]byte(username))
	passwordHash := sha256.Sum256([]byte(password))

	usernameMatch := subtle.ConstantTimeCompare(expectedUsernameHash[:], usernameHash[:]) == 1
	passwordMatch := subtle.ConstantTimeCompare(expectedPasswordHash[:], passwordHash[:]) == 1

	return usernameMatch && passwordMatch
}
