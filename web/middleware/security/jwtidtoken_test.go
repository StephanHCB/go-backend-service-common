package security

import (
	"github.com/StephanHCB/go-backend-service-common/docs"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// private key was thrown away - if you need a new key for changes to the test cases, just have jwt.io roll you one

const tstJwtPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ
0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg
cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc
mwIDAQAB
-----END PUBLIC KEY-----`

// expired 2019
const tstTokenBodyExpired = `{
  "iss": "myissuer",
  "sub": "1234567890",
  "exp": 1516239023,
  "name": "John Doe",
  "email": "john.doe@example.com",
  "groups": ["somegroup"],
  "iat": 1516239022
}`

// expires 2049
const tstTokenBodyCurrent = `{
  "iss": "myissuer",
  "sub": "1234567890",
  "exp": 2516239023,
  "name": "John Doe",
  "email": "john.doe@example.com",
  "groups": ["somegroup"],
  "iat": 1516239022
}`

const tstJwtValidToken = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJteWlzc3VlciIsInN1YiI6IjEyMzQ1Njc4OTAiLCJleHAiOjI1MTYyMzkwMjMsIm5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJncm91cHMiOlsic29tZWdyb3VwIl0sImlhdCI6MTUxNjIzOTAyMn0.TkJR5ik714nFfgP4K40RgsOaziGIabTcQ_GOBYKJi53GnGv9Obn9ROIqIFiSiYS_TWYagRxK3FkW9pLeTME3lx064eOF7WLi6CbSQCpVghM1oJKVdwXoqksT6B3YwpPdm2GhWdQ-aGMukjadjbigFNZnjAjOqKNGgoYqz02BD25KLfWZIFN9MLeBTJj5SiFy1PorruuxPCLcIHg-HPczWeS9ux_W8yaQtgfPHvwMUpW4e0sPlO7ipJwQqIXMjwsCfvdnQODViGxkPaDwqH80nkiv9bd72M7OwM2O4He1Z1kaED1PtISNUhKGSvAhjDk8yNOVRZFeNoQUemTAb09eMA`
const tstJwtExpiredToken = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJteWlzc3VlciIsInN1YiI6IjEyMzQ1Njc4OTAiLCJleHAiOjE1MTYyMzkwMjMsIm5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJncm91cHMiOlsic29tZWdyb3VwIl0sImlhdCI6MTUxNjIzOTAyMn0.mFzOoIfHADN6DglrQCuYsGJJ6Mx6r12XoqpymJ0KL1jLC9mVuK0yGQ2yUpoV9YeaKOd7A2_WHZZTohGa2QavRimGe04xDHcD7sUeo-WemQ85sWtfb0XtAuJJKMUb6qu28LCYKU1x5lZGyqHiPJmRqGqsbTXZjPMM7e2gdPqKcVuYSwHQzArZh8DFYu9Cgx8j1V_mrcWcuChXZfwpgMqGd5bwNb-t90b9pxF2vmqhq6WNnnyfV-MB-XmJ2geLyw0rPgaNrZtgFYAkepf3Qr4q9edlNntAAFPuehEVgPdz-WrWx3Iux5OFAa97L3S4kzioWq4ZoRRONNv4ICyMBCFteA`
const tstJwtInvalidB64Token = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyXxXxXHEREXxXxXJpc3MiOiJteWlzc3VlciIsInN1YiI6IjEyMzQ1Njc4OTAiLCJleHAiOjI1MTYyMzkwMjMsIm5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJncm91cHMiOlsic29tZWdyb3VwIl0sImlhdCI6MTUxNjIzOTAyMn0.TkJR5ik714nFfgP4K40RgsOaziGIabTcQ_GOBYKJi53GnGv9Obn9ROIqIFiSiYS_TWYagRxK3FkW9pLeTME3lx064eOF7WLi6CbSQCpVghM1oJKVdwXoqksT6B3YwpPdm2GhWdQ-aGMukjadjbigFNZnjAjOqKNGgoYqz02BD25KLfWZIFN9MLeBTJj5SiFy1PorruuxPCLcIHg-HPczWeS9ux_W8yaQtgfPHvwMUpW4e0sPlO7ipJwQqIXMjwsCfvdnQODViGxkPaDwqH80nkiv9bd72M7OwM2O4He1Z1kaED1PtISNUhKGSvAhjDk8yNOVRZFeNoQUemTAb09eMA`
const tstJwtInvalidSignatureToken = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJteWlzc3VlciIsInN1YiI6IjEyMzQ1Njc4OTAiLCJleHAiOjI1MTYyMzkwMjMsIm5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJncm91cHMiOlsic29tZWdyb3VwIl0sImlhdCI6MTUxNjIzOTAyMn0.TkJR5ik714nFfgP4K40RgsOaziGIabTcQ_GOBYKJi53GnGv9Obn9ROIqIFiSiYS_TWYagRxK3FkW9pLeTME3lx064eOF7WLi6CbSQCpVghM1oJKVdwXoqksT6B3YwpPdm2GhWdQ-aGMukjadjbigFNZnjAjOqKNGgoYqz02BD25KLfWZIFN9MLeBTJj5SiFy1PorruuxPCLcIHg-HPczWeS9ux_W8yaQtgfPHvwMUpW4e0sPlO7ipJwQqIXMjwsCfvdnQODViGxkPaDwqH80nkiv9bd72M7OwM2O4He1Z1kaED1PtISNUhKGSvAhjDk8yNOVRZFeNoQUemTAb09`
const tstJwtInvalidJsonToken = `Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJteWlzc3VlciIsInN1YiI6IjEyMzQ1Njc4OTAiLCJleHAiOjI1MTYyMzkwMjMsIm5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJncm91cHMiOlsic29tZWdyb3VwIl0sImlhdCI6MTUxNjIzOTAyM.C6b1wYFrrBfm_1fn1akK_dXEPXoy7RxBTvzZ3zIJxoJ5BdXkhef-vPpiGewVuqakF5VDrhqWdciIHgnG0eHFgBJcmypwD1TNMaDP238M1TjAsM36ZikHCONQMQ6AUC3emrCBtx5eRKvwGu8q_8iaBW3D-4n_vwIDPAGuSFiR1oIKTKPLCETQlZn-fCA_EKjYiab9oI8BJwEVHwPidBlMD81ImFvlghm8noeEjVxUpwB39dKZ6yj8T5KMCoqt-DX1oZqenf2f7ZdJT3qsopv3iEnRZN0Zxu0oMxOAHJmp1ITuKEWeMBfYp41NmZsluIvgboMRq9skjk`

const tstOtherAuthType = `Bringer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJteWlzc3VlciIsInN1YiI6IjEyMzQ1Njc4OTAiLCJleHAiOjI1MTYyMzkwMjMsIm5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJncm91cHMiOlsic29tZWdyb3VwIl0sImlhdCI6MTUxNjIzOTAyMn0.TkJR5ik714nFfgP4K40RgsOaziGIabTcQ_GOBYKJi53GnGv9Obn9ROIqIFiSiYS_TWYagRxK3FkW9pLeTME3lx064eOF7WLi6CbSQCpVghM1oJKVdwXoqksT6B3YwpPdm2GhWdQ-aGMukjadjbigFNZnjAjOqKNGgoYqz02BD25KLfWZIFN9MLeBTJj5SiFy1PorruuxPCLcIHg-HPczWeS9ux_W8yaQtgfPHvwMUpW4e0sPlO7ipJwQqIXMjwsCfvdnQODViGxkPaDwqH80nkiv9bd72M7OwM2O4He1Z1kaED1PtISNUhKGSvAhjDk8yNOVRZFeNoQUemTAb09eMA`

func TestJwtIdTokenValidatorMiddleware_NoAuth_FallThrough(t *testing.T) {
	docs.Description("jwt middleware lets requests through with missing authorization header")
	tstJwtIdTokenTestcase(t, noAuth, true, false)
}

func TestJwtIdTokenValidatorMiddleware_OtherAuthType_FallThrough(t *testing.T) {
	docs.Description("jwt middleware lets requests through with authorization header of non-Bearer type")
	tstJwtIdTokenTestcase(t, tstOtherAuthType, true, false)
}

func TestJwtIdTokenValidatorMiddleware_ValidToken_Ok(t *testing.T) {
	docs.Description("jwt middleware successfully authorizes requests with valid token")
	tstJwtIdTokenTestcase(t, tstJwtValidToken, true, true)
}

func TestJwtIdTokenValidatorMiddleware_ExpiredToken_Reject(t *testing.T) {
	docs.Description("jwt middleware rejects requests with expired token")
	tstJwtIdTokenTestcase(t, tstJwtExpiredToken, false, false)
}

func TestJwtIdTokenValidatorMiddleware_B64ErrorToken_Reject(t *testing.T) {
	docs.Description("jwt middleware rejects requests with B64 malformed token")
	tstJwtIdTokenTestcase(t, tstJwtInvalidB64Token, false, false)
}

func TestJwtIdTokenValidatorMiddleware_InvalidSignatureToken_Reject(t *testing.T) {
	docs.Description("jwt middleware rejects requests with invalid signature token")
	tstJwtIdTokenTestcase(t, tstJwtInvalidSignatureToken, false, false)
}

func TestJwtIdTokenValidatorMiddleware_JsonErrorToken_Reject(t *testing.T) {
	docs.Description("jwt middleware rejects requests with json malformed token")
	tstJwtIdTokenTestcase(t, tstJwtInvalidJsonToken, false, false)
}

// --- helpers ---

func tstJwtIdTokenTestcase(t *testing.T, authorization string, shouldGoThrough bool, shouldAuthorize bool) {
	wentThrough := false
	cut := tstConstructJwtIdTokenHandlerUnderTest(t, shouldAuthorize, &wentThrough)

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

func tstConstructJwtIdTokenHandlerUnderTest(t *testing.T, expectAuthorized bool, wentThrough *bool) http.Handler {
	verifyingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedClaims := CustomClaims{
			Name:   "John Doe",
			Email:  "john.doe@example.com",
			Groups: []string{"somegroup"},
		}

		ctx := r.Context()

		claims := GetClaims(ctx)

		if wentThrough != nil {
			// signal to the test that the middleware did not block the request
			*wentThrough = true
		}

		if expectAuthorized {
			require.NotNil(t, claims)
			require.EqualValues(t, expectedClaims, claims.CustomClaims)
		} else {
			require.Nil(t, claims)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	parsedKeys, err := ParsePublicKeysFromPEM([]string{tstJwtPublicKey})
	require.Nil(t, err)

	options := JwtIdTokenValidatorMiddlewareOptions{
		PublicKeys: parsedKeys,
	}
	middlewareUnderTest := JwtIdTokenValidatorMiddleware(options)

	return middlewareUnderTest(verifyingHandler)
}
