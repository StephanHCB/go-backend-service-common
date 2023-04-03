package util

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-backend-service-common/api/apierrors"
	"net/http"
	"time"
)

func UnauthorizedErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, logMessage string, timeStamp time.Time) {
	aulogging.Logger.Ctx(ctx).Info().Printf("unauthorized: %s", logMessage)
	apierrors.ErrorHandler(ctx, w, r, http.StatusUnauthorized, "unauthorized", "missing or invalid Authorization header (JWT bearer token expected) or token invalid or expired", nil, timeStamp)
}
