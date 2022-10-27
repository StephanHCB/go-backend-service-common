package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"net/http"
)

// Key to use when setting the request ID.
type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0

// RequestIDHeader is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var RequestIDHeader = "X-Request-Id"

// NewRequestIDFunc allows overriding the generator function for a new request id.
var NewRequestIDFunc = NewRequestID

// RequestID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = NewRequestIDFunc()
		}
		ctx = PutReqID(ctx, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// GetReqID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// PutReqID puts the given requestID in the context under the correct key.
// Returns a child context with the value set
func PutReqID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// NewRequestID generates a fresh request ID.
func NewRequestID() string {
	requestID := aulogging.DefaultRequestIdValue

	var buf [4]byte
	_, err := rand.Read(buf[:])
	if err == nil {
		requestID = hex.EncodeToString(buf[:])
	}

	return requestID
}
