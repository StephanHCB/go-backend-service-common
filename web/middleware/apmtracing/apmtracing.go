package apmtracing

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/StephanHCB/go-autumn-logging-zerolog/loggermiddleware"
	auapmlogging "github.com/StephanHCB/go-autumn-restclient-apm/implementation/logging"
	auapmmiddleware "github.com/StephanHCB/go-autumn-restclient-apm/implementation/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.elastic.co/apm/module/apmchiv5/v2"
	"go.elastic.co/apm/v2"
	"net/http"
)

type ApmMiddlewareOptions struct {
	ElasticApmEnabled bool
	PlainLogging      bool
}

func BuildApmMiddleware(ctx context.Context, options ApmMiddlewareOptions) func(http.Handler) http.Handler {
	// add apm middleware, because we rely on having a trace context in the context for trace logging and trace propagation to work.
	if !options.ElasticApmEnabled {
		// if apm is not configured, we use a discardTracer that does not send any traces
		err := auapmmiddleware.SetupDiscardTracer()
		if err != nil {
			aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Print("setting up discard tracer failed - continuing with default tracer: %s", err.Error())
		} else {
			aulogging.Logger.Ctx(ctx).Info().Print("successfully set up discard tracer because Elastic APM is not configured")
		}
	}

	return GetApmMiddlewareAfterSetup(ctx, options)
}

func GetApmMiddlewareAfterSetup(_ context.Context, _ ApmMiddlewareOptions) func(http.Handler) http.Handler {
	return apmchiv5.Middleware()
}

func ConfigureContextLoggingForApm(ctx context.Context, options ApmMiddlewareOptions) loggermiddleware.AddZerologLoggerToContextOptions {
	if options.PlainLogging {
		// override requestIdRetriever to see APM trace ids in plain logging
		// (for json logging, we set up additional fields instead, so you also have the standard Zipkin request id)
		aulogging.RequestIdRetriever = auapmlogging.ExtractTraceId

		// no additional fields for plain logging (readability)
		return loggermiddleware.AddZerologLoggerToContextOptions{}
	} else {
		return loggermiddleware.AddZerologLoggerToContextOptions{
			CustomJsonLogFields: []loggermiddleware.CustomJsonLogField{
				customJsonLogField(auapmlogging.TraceIdLogFieldName, auapmlogging.ExtractTraceId),
				customJsonLogField(auapmlogging.TransactionIdLogFieldName, auapmlogging.ExtractTransactionId),
				customJsonLogField(auapmlogging.SpanIdLogFieldName, auapmlogging.ExtractSpanId),
			},
		}
	}
}

func customJsonLogField(name string, extractor func(context.Context) string) loggermiddleware.CustomJsonLogField {
	return loggermiddleware.CustomJsonLogField{
		LogFieldName: name,
		ValueExtractor: func(r *http.Request) string {
			return extractor(r.Context())
		},
	}
}

// StartTransaction starts a completely new APM transaction and places it in the provided context.
// Also places an APM enabled logger in the context, based on existing logger, if present.
//
// cancelFunc must be called via defer to close the transaction!
func StartTransaction(ctx context.Context, name string, transactionType string) (context.Context, context.CancelFunc) {
	tx := apm.DefaultTracer().StartTransaction(name, transactionType)
	apmCtx := apm.ContextWithTransaction(ctx, tx)

	sourceLogger := zerolog.Ctx(ctx)
	if sourceLogger == nil {
		sourceLogger = &log.Logger
	}

	traceLogger := addTransactionLoggerFields(apmCtx, *sourceLogger)
	traceCtx := traceLogger.WithContext(apmCtx)

	return traceCtx, func() {
		tx.End()
	}
}

// StartSpan starts a span for an already running APM transaction.
// Also places an APM enabled logger in the context, based on existing logger, if present.
func StartSpan(ctx context.Context, name string, spanType string) (context.Context, context.CancelFunc) {
	span, spanCtx := apm.StartSpan(ctx, name, spanType)

	sourceLogger := zerolog.Ctx(ctx)
	if sourceLogger == nil {
		sourceLogger = &log.Logger
	}

	traceLogger := addSpanLoggerFields(spanCtx, *sourceLogger)
	traceCtx := traceLogger.WithContext(spanCtx)

	return traceCtx, func() {
		span.End()
	}
}

func addTransactionLoggerFields(ctx context.Context, sourceLogger zerolog.Logger) zerolog.Logger {
	builder := sourceLogger.With()
	if auzerolog.IsJson {
		if aulogging.RequestIdRetriever != nil {
			requestId := aulogging.RequestIdRetriever(ctx)
			builder = builder.Str(loggermiddleware.RequestIdFieldName, requestId)
		}
		// TODO log name in a suitable field? Might be useful!
		builder = builder.Str(auapmlogging.TraceIdLogFieldName, auapmlogging.ExtractTraceId(ctx))
		builder = builder.Str(auapmlogging.TransactionIdLogFieldName, auapmlogging.ExtractTransactionId(ctx))
	} else {
		// use request id field for trace id because plain logger does not support custom fields
		builder = builder.Str(loggermiddleware.RequestIdFieldName, auapmlogging.ExtractTraceId(ctx))
	}
	sublogger := builder.Logger()
	return sublogger
}

func addSpanLoggerFields(ctx context.Context, sourceLogger zerolog.Logger) zerolog.Logger {
	builder := sourceLogger.With()
	if auzerolog.IsJson {
		builder = builder.Str(auapmlogging.SpanIdLogFieldName, auapmlogging.ExtractSpanId(ctx))
	}
	logger := builder.Logger()
	return logger
}
