package requestmetrics

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
	"time"
)

var (
	RequestCounterName  = "http_server_requests_seconds_count"
	RequestDurationName = "http_server_requests_seconds_sum"

	reqs    *prometheus.CounterVec
	latency *prometheus.SummaryVec
)

func Setup() {
	reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: RequestCounterName,
			Help: "Number of incoming HTTP requests processed, partitioned by status code, method and HTTP path (grouped by patterns).",
		},
		[]string{"method", "outcome", "status", "uri"},
	)
	prometheus.MustRegister(reqs)

	latency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: RequestDurationName,
			Help: "How long it took to process requests, partitioned by status code, method and HTTP path (grouped by patterns).",
		},
		[]string{"method", "outcome", "status", "uri"},
	)
	prometheus.MustRegister(latency)
}

func RecordRequestMetrics(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.Replace(routePattern, "/*/", "/", -1)

		reqs.WithLabelValues(r.Method, outcome(ww.Status()), fmt.Sprintf("%d", ww.Status()), routePattern).Inc()
		latency.WithLabelValues(r.Method, outcome(ww.Status()), fmt.Sprintf("%d", ww.Status()), routePattern).Observe(float64(time.Since(start).Microseconds()) / 1000000)
	}
	return http.HandlerFunc(fn)
}

func outcome(status int) string {
	if status < 400 {
		return "SUCCESS"
	} else if status < 500 {
		return "CLIENT_ERROR"
	} else {
		return "SERVER_ERROR"
	}
}
