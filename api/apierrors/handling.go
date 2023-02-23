package apierrors

import (
	"context"
	"encoding/json"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-backend-service-common/web/util/media"
	"github.com/go-http-utils/headers"
	"net/http"
)

// HandleError is a common error handler for all errors declared in this package.
//
// We make you pass in a list of expectedType checks so each handler function documents what api errors it
// expects to happen, any unexpected type will cause a 500 until you update the code (and don't forget
// the OpenAPI spec while you're at it)!
//
// Pass any number of the IsXyz functions from this package for expectedTypes.
func HandleError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error, expectedTypes ...func(error) bool) {
	annotatedError, ok := err.(AnnotatedError)
	if ok {
		for _, typeCheck := range expectedTypes {
			if typeCheck(err) {
				msg := annotatedError.ApiError().Message
				errorHandler(ctx, w, r, *annotatedError.ApiError().ErrorCode, *msg)
				return
			}
		}
	}
	// ensure 500 if a handler throws a type of error not documented in the OpenAPI spec
	unexpectedErrorHandler(ctx, w, r, err)
}

func unexpectedErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	aulogging.Logger.Ctx(ctx).Error().WithErr(err).Printf("unexpected error")
	errorHandler(ctx, w, r, http.StatusInternalServerError, err.Error())
}

func errorHandler(ctx context.Context, w http.ResponseWriter, _ *http.Request, status int32, details string) {
	response := ApiErrorV2{
		ErrorCode: &status,
		Message:   &details,
	}

	w.Header().Set(headers.ContentType, media.ContentTypeApplicationJson)
	w.WriteHeader(int(status))
	writeJson(ctx, w, response)
}

func writeJson(ctx context.Context, w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error while encoding json response: %v", err)
		// can't change status anymore, in the middle of the response now
	}
}
