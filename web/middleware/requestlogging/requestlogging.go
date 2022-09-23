// overriding parts of chi's middleware.Logger because we want to use go-autumn-logging
package requestlogging

import (
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	auloggingapi "github.com/StephanHCB/go-autumn-logging/api"
	"github.com/StephanHCB/go-backend-service-common/web/middleware/requestid"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func Setup() {
	middleware.DefaultLogger = middleware.RequestLogger(&zerologLogFormatter{})
}

// --- implement middleware.LogFormatter

type zerologLogFormatter struct {
}

const (
	StatusCodeFieldName            = "http.response.status_code"
	UserAgentFieldName             = "user_agent.original"
	ResponseLatencyMicrosFieldName = "event.duration"
	LoggerNameFieldName            = "log.logger"
)

func (l *zerologLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &zerologLogEntry{
		zerologLogFormatter: l,
		request:             r,
	}
	entry.requestId = requestid.GetReqID(r.Context())
	entry.method = r.Method
	entry.path = r.URL.Path
	entry.ip = r.RemoteAddr
	entry.userAgent = r.UserAgent()

	return entry
}

// --- implement middleware.LogEntry

type zerologLogEntry struct {
	*zerologLogFormatter
	request   *http.Request
	requestId string
	method    string
	path      string
	ip        string
	userAgent string
}

func (l *zerologLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	msg := "request"

	ctxLogger := aulogging.Logger.Ctx(l.request.Context())
	var e auloggingapi.LeveledLoggingImplementation
	switch {
	case status >= http.StatusInternalServerError:
		e = ctxLogger.Error()
	default:
		e = ctxLogger.Info()
	}

	if auzerolog.IsJson {
		e.With(StatusCodeFieldName, fmt.Sprintf("%d", status)).
			With(ResponseLatencyMicrosFieldName, fmt.Sprintf("%d", elapsed.Microseconds())).
			With(UserAgentFieldName, l.userAgent).
			With(LoggerNameFieldName, "request.incoming").
			Print(msg)
	} else {
		// console friendly version
		msg = fmt.Sprintf("request %s %s -> %d (%d μs)", l.method, l.path, status, elapsed.Microseconds())
		e.Print(msg)
	}
}

func (l *zerologLogEntry) Panic(v interface{}, stack []byte) {
	panicEntry := l.NewLogEntry(l.request).(*zerologLogEntry)

	msg := "Request Panic"

	e := aulogging.Logger.NoCtx().Panic()

	e.With(UserAgentFieldName, panicEntry.userAgent).
		Print(msg)
}
