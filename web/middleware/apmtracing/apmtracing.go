package apmtracing

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-autumn-logging-zerolog/loggermiddleware"
	auapmlogging "github.com/StephanHCB/go-autumn-restclient-apm/implementation/logging"
	auapmmiddleware "github.com/StephanHCB/go-autumn-restclient-apm/implementation/middleware"
	"go.elastic.co/apm/module/apmchiv5/v2"
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
