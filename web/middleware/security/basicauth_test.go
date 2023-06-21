package security

import (
	"github.com/StephanHCB/go-backend-service-common/docs"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const noAuth = ""
const basicAuthValid = "Basic dGVzdHVzZXI6dGVzdHB3"           // echo -n "testuser:testpw" | base64
const basicAuthInvalidCred = "Basic dGVzdHVzZXI6d3Jvbmdwdw==" // echo -n "testuser:wrongpw" | base64
const basicAuthInvalidEncoding = "Basic dGVzdHVzZXI6dGVzdqr"
const otherAuth = "Bearer something"

func TestBasicAuthValidatorMiddleware_NoAuth_FallThrough(t *testing.T) {
	docs.Description("basic auth middleware lets requests through with missing authorization header")
	tstBasicAuthTestcase(t, noAuth, true, false)
}

func TestBasicAuthValidatorMiddleware_ValidAuth_Authorizes(t *testing.T) {
	docs.Description("basic auth middleware allows valid basic auth through")
	tstBasicAuthTestcase(t, basicAuthValid, true, true)
}

func TestBasicAuthValidatorMiddleware_InvalidCred_Fails(t *testing.T) {
	docs.Description("basic auth middleware refuses wrong credentials")
	tstBasicAuthTestcase(t, basicAuthInvalidCred, false, false)
}

func TestBasicAuthValidatorMiddleware_InvalidEncoding_Fails(t *testing.T) {
	docs.Description("basic auth middleware refuses invalid base64 encoding")
	tstBasicAuthTestcase(t, basicAuthInvalidEncoding, false, false)
}

func TestBasicAuthValidatorMiddleware_OtherAuth_FallThrough(t *testing.T) {
	docs.Description("basic auth middleware lets other auth methods fall through")
	tstBasicAuthTestcase(t, otherAuth, true, false)
}

// --- helpers ---

func tstBasicAuthTestcase(t *testing.T, authorization string, shouldGoThrough bool, shouldAuthorize bool) {
	wentThrough := false
	cut := tstConstructBasicAuthHandlerUnderTest(t, shouldAuthorize, &wentThrough)

	r := httptest.NewRequest(http.MethodGet, "/v1/api", nil)
	if authorization != "" {
		r.Header.Add(headers.Authorization, authorization)
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

func tstConstructBasicAuthHandlerUnderTest(t *testing.T, expectAuthorized bool, wentThrough *bool) http.Handler {
	specifiedClaims := CustomClaims{
		Name:   "testuser",
		Email:  "testuser@example.com",
		Groups: []string{"testgroup"},
	}

	verifyingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wentThrough != nil {
			// signal to the test that the middleware did not block the request
			*wentThrough = true
		}

		ctx := r.Context()

		claims := GetClaims(ctx)
		if expectAuthorized {
			require.NotNil(t, claims)
			require.EqualValues(t, specifiedClaims, claims.CustomClaims)
		} else {
			require.Nil(t, claims)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	options := BasicAuthMiddlewareOptions{
		BasicAuthUsername: "testuser",
		BasicAuthPassword: "testpw",
		BasicAuthClaims:   specifiedClaims,
	}
	middlewareUnderTest := BasicAuthValidatorMiddleware(options)

	return middlewareUnderTest(verifyingHandler)
}
