package corsheader

import (
	"github.com/StephanHCB/go-backend-service-common/web/middleware/requestid"
	"github.com/go-http-utils/headers"
	"net/http"
	"strings"
)

func CorsHandling(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headers.AccessControlAllowOrigin, "*")

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
