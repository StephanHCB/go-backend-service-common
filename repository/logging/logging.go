package logging

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	auloggingapi "github.com/StephanHCB/go-autumn-logging/api"
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
	"github.com/StephanHCB/go-backend-service-common/web/middleware/requestid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

type LoggingImpl struct {
	Configuration repository.Configuration
	Metrics       *prometheus.CounterVec
}

var LogCounterName = "logging_events_total"

func (l *LoggingImpl) Setup() {
	aulogging.RequestIdRetriever = requestid.GetReqID

	l.Metrics = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: LogCounterName,
			Help: "How many log entries were written per level.",
		},
		[]string{"level"},
	)
	prometheus.MustRegister(l.Metrics)

	aulogging.LogEventCallback = l.loggingCallback

	if l.Configuration.PlainLogging() {
		aulogging.DefaultRequestIdValue = "00000000"
		auzerolog.RequestIdFieldName = "trace.id"
		auzerolog.SetupPlaintextLogging()
		aulogging.Logger.NoCtx().Info().Print("switching to developer friendly console log")
	} else {
		// stay with JSON logging and add ECS service.id field
		l.CustomSetupJsonLogging(l.Configuration.ApplicationName())
	}
}

func (l *LoggingImpl) loggingCallback(_ context.Context, level string, _ string, _ error, _ map[string]string) {
	l.Metrics.WithLabelValues(strings.ToLower(level)).Inc()
}

// override auzerolog.SetupJsonLogging so we can get as close to the other services as possible

func (l *LoggingImpl) CustomSetupJsonLogging(serviceName string) {
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.LevelFieldName = "log.level"
	zerolog.MessageFieldName = "message" // correct by default

	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service.name", serviceName).
		Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.LevelTraceValue = "TRACE"
	zerolog.LevelDebugValue = "DEBUG"
	zerolog.LevelInfoValue = "INFO"
	zerolog.LevelWarnValue = "WARN"
	zerolog.LevelErrorValue = "ERROR"
	zerolog.LevelFatalValue = "FATAL"
	zerolog.LevelPanicValue = "FATAL"

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z"

	auzerolog.IsJson = true
}

// alternative Setup function for testing that records log entries instead of writing them to console
func (l *LoggingImpl) SetupForTesting() {
	auzerolog.SetupForTesting()
}

func (l *LoggingImpl) Logger() auloggingapi.LoggingImplementation {
	return aulogging.Logger
}
