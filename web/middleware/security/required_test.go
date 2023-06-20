package security

import (
	"github.com/StephanHCB/go-backend-service-common/docs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRequiredMiddleware_Auth_Allow(t *testing.T) {
	docs.Description("required auth middleware lets authorized requests through")
	tstAuthRequiredTestcase(t, http.MethodPut, "/v2/irrelevant", true, true)
}

func TestAuthRequiredMiddleware_NoAuth_BlockedMethod(t *testing.T) {
	docs.Description("required auth middleware blocks requests not on allow list")
	tstAuthRequiredTestcase(t, http.MethodPost, "/v2/matches/glob", false, false)
}

func TestAuthRequiredMiddleware_NoAuth_BlockedPath(t *testing.T) {
	docs.Description("required auth middleware blocks requests not on allow list (path glob)")
	tstAuthRequiredTestcase(t, http.MethodGet, "/v17/does/not/match", false, false)
}

func TestAuthRequiredMiddleware_NoAuth_OnAllowlist(t *testing.T) {
	docs.Description("required auth middleware lets unauthorized requests through if matching allow list")
	tstAuthRequiredTestcase(t, http.MethodGet, "/v2/matches/glob", false, true)
}

// --- helpers ---

func tstAuthRequiredTestcase(t *testing.T, method string, target string, authorized bool, shouldGoThrough bool) {
	wentThrough := false
	cut := tstConstructRequiredHandlerUnderTest(t, &wentThrough)

	r := httptest.NewRequest(method, target, nil)
	if authorized {
		specifiedClaims := AllClaims{
			RegisteredClaims: jwt.RegisteredClaims{},
			CustomClaims: CustomClaims{
				Name:  "someone",
				Email: "someone@example.com",
			},
		}
		ctx := PutClaims(r.Context(), &specifiedClaims)
		r = r.WithContext(ctx)
	}
	w := httptest.NewRecorder()
	cut.ServeHTTP(w, r)

	require.Equal(t, shouldGoThrough, wentThrough)
	if shouldGoThrough {
		require.Equal(t, http.StatusNoContent, w.Code)
	} else {
		require.Equal(t, http.StatusUnauthorized, w.Code)
	}
}

func tstConstructRequiredHandlerUnderTest(t *testing.T, wentThrough *bool) http.Handler {
	verifyingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wentThrough != nil {
			// signal to the test that the middleware did not block the request
			*wentThrough = true
		}
		w.WriteHeader(http.StatusNoContent)
	})

	options := AuthRequiredMiddlewareOptions{
		AllowUnauthorized: []string{"GET /v2/.*", "POST /v17/.*/allowed"},
	}
	middlewareUnderTest := AuthRequiredMiddleware(options)

	return middlewareUnderTest(verifyingHandler)
}
