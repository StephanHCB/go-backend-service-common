package requestidinresponse

import (
	"github.com/StephanHCB/go-backend-service-common/web/middleware/requestid"
	"net/http"
)

func AddRequestIdHeaderToResponse(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestId := requestid.GetReqID(ctx)
		if requestId != "" {
			w.Header().Set(requestid.RequestIDHeader, requestId)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
