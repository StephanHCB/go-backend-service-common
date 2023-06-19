package security

import (
	"context"
	"encoding/json"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-backend-service-common/api"
	"github.com/StephanHCB/go-backend-service-common/api/apierrors"
	"github.com/StephanHCB/go-backend-service-common/web/util/media"
	"github.com/go-http-utils/headers"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type ctxSecurityKeyType int

const (
	RawTokenKey ctxSecurityKeyType = 0
	ClaimsKey   ctxSecurityKeyType = 1
)

type CustomClaims struct {
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Groups []string `json:"groups"`
}

type AllClaims struct {
	jwt.RegisteredClaims
	CustomClaims
}

// Now exported for testing
var Now = time.Now

// GetRawToken returns the raw token from the given context if one is present.
//
// Returns the empty string if the context contains no valid token.
func GetRawToken(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if token, ok := ctx.Value(RawTokenKey).(string); ok {
		return token
	}
	return ""
}

// PutRawToken places a raw token in the context under the correct key.
//
// Returns a child context with the token set.
//
// Exposed for testing only.
func PutRawToken(ctx context.Context, rawToken string) context.Context {
	return context.WithValue(ctx, RawTokenKey, rawToken)
}

// GetClaims returns the raw token from the given context if one is present.
//
// Returns the empty string if the context contains no valid token.
func GetClaims(ctx context.Context) *AllClaims {
	if ctx == nil {
		return nil
	}
	if claimsPtr, ok := ctx.Value(ClaimsKey).(*AllClaims); ok {
		return claimsPtr
	}
	return nil
}

// PutClaims places a raw token in the context under the correct key.
//
// Returns a child context with the token set.
//
// Exposed for testing only.
func PutClaims(ctx context.Context, claimsPtr *AllClaims) context.Context {
	return context.WithValue(ctx, ClaimsKey, claimsPtr)
}

func IsAuthenticated(ctx context.Context, logMessage string, timestamp time.Time) apierrors.AnnotatedError {
	claimsPtr := GetClaims(ctx)
	if claimsPtr == nil {
		aulogging.Logger.Ctx(ctx).Info().Printf("unauthorized: %s", logMessage)
		return apierrors.NewUnauthorisedError("unauthorized", "missing or invalid Authorization header (JWT bearer token expected) or token invalid or expired", nil, timestamp)
	}
	return nil
}

func HasGroup(ctx context.Context, group string, logMessage string, timestamp time.Time) apierrors.AnnotatedError {
	err := apierrors.NewForbiddenError("forbidden", "you are not authorized for this operation", nil, timestamp)
	if group == "" {
		return nil
	}
	claimsPtr := GetClaims(ctx)
	if claimsPtr == nil {
		aulogging.Logger.Ctx(ctx).Info().Printf("forbidden: %s", logMessage)
		return err
	}
	if !contains(claimsPtr.Groups, group) {
		return err
	}
	return nil
}

func Name(ctx context.Context) string {
	claimsPtr := GetClaims(ctx)
	if claimsPtr == nil {
		return ""
	}
	return claimsPtr.Name
}

func Email(ctx context.Context) string {
	claimsPtr := GetClaims(ctx)
	if claimsPtr == nil {
		return ""
	}
	return claimsPtr.Email
}

func Subject(ctx context.Context) string {
	claimsPtr := GetClaims(ctx)
	if claimsPtr == nil {
		return ""
	}
	return claimsPtr.RegisteredClaims.Subject
}

func contains(haystack []string, needle string) bool {
	if len(haystack) == 0 {
		return false
	}
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

// error handlers

func unauthorizedErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, logMessage string, timeStamp time.Time) {
	aulogging.Logger.Ctx(ctx).Info().Printf("unauthorized: %s", logMessage)
	errorHandler(ctx, w, r, "unauthorized", http.StatusUnauthorized, "missing or invalid Authorization header (JWT bearer token expected) or token invalid or expired", timeStamp)
}

func errorHandler(ctx context.Context, w http.ResponseWriter, _ *http.Request, msg string, status int, details string, timestamp time.Time) {
	detailsPtr := &details
	if details == "" {
		detailsPtr = nil
	}
	response := &api.ErrorDto{
		Details:   detailsPtr,
		Message:   &msg,
		Timestamp: &timestamp,
	}
	w.Header().Set(headers.ContentType, media.ContentTypeApplicationJson)
	w.WriteHeader(status)
	WriteJson(ctx, w, response)
}

func WriteJson(ctx context.Context, w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error while encoding json response: %v", err)
		// can't change status anymore, in the middle of the response now
	}
}
