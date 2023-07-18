package contexthelper

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auapmlogging "github.com/StephanHCB/go-autumn-restclient-apm/implementation/logging"
	"github.com/StephanHCB/go-backend-service-common/web/middleware/apmtracing"
	"github.com/StephanHCB/go-backend-service-common/web/middleware/requestid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.elastic.co/apm/v2"
)

// regarding APM tracing, see here:
// https://github.com/StephanHCB/go-autumn-restclient-apm#example-integration-into-go-autumn-chi-service

// StandaloneContext creates a fully configured context for processing that is not based on an existing context.
//
// A context is built from scratch, a child of context.Background(). Logger and request id are created and provided.
// We also start a completely new tracing transaction.
//
// name is used as the name for the trace transaction and in the logger.
//
// traceTransactionType is used as the trace transaction type. Typical values are "scheduled", "request",
// "backgroundJob". See https://www.elastic.co/guide/en/apm/guide/current/data-model-transactions.html for details.
//
// Usage example:
//
//	ctx, cancel := contexthelper.StandaloneContext("RefreshMetadata", "scheduled")
//	defer cancel()
//
// The returned cancel function will close the transaction, clean up, then finally cancel the context.
func StandaloneContext(name string, traceTransactionType string) (ctx context.Context, cancel context.CancelFunc) {
	ctx = context.Background()

	// child context so we do not cancel global context
	var fullCancel context.CancelFunc
	ctx, fullCancel = context.WithCancel(ctx)

	// construct a new Zipkin request id and put in context
	newRequestId := requestid.NewRequestIDFunc()
	ctx = requestid.PutReqID(ctx, newRequestId)

	// start APM transaction, so we get a trace id, and add APM enabled logger to the context
	var apmCancel context.CancelFunc
	ctx, apmCancel = apmtracing.StartTransaction(ctx, name, traceTransactionType)

	cancel = func() {
		apmCancel()
		fullCancel()
	}

	return ctx, cancel
}

// AsyncCopyRequestContext creates a fully configured context for asynchronous processing that is started
// by a web request, or from a previously generated standalone context.
//
// The context is built as a child of context.Background(), using information from sourceCtx.
// This means, cancellation of sourceCtx will not cancel the returned context. Logger and request id are copied over.
// We also start a new tracing transaction, logging a line to link it to the old transaction.
//
// Note: we cannot just create a span under the original transaction, because it will end when the original
// request completes, which would end all spans.
//
// name is used as the name for the trace transaction.
//
// traceTransactionType is used as the trace transaction type. Typical values are "scheduled", "request",
// "backgroundJob". See https://www.elastic.co/guide/en/apm/guide/current/data-model-transactions.html for details.
//
// Usage example:
//
//	ctx, cancel := contexthelper.AsyncCopyRequestContext(r.Context(), "Async Webhook Reply 1", "request")
//	defer cancel()
//
// The returned cancel function will close the span, clean up, then finally cancel the context.
func AsyncCopyRequestContext(sourceCtx context.Context, name string, traceTransactionType string) (ctx context.Context, cancel context.CancelFunc) {
	ctx = context.Background()

	// child context so we do not cancel global context
	var fullCancel context.CancelFunc
	ctx, fullCancel = context.WithCancel(ctx)

	// carry over Zipkin request id, if any
	requestId := requestid.GetReqID(sourceCtx)
	if requestId != "" {
		ctx = requestid.PutReqID(ctx, requestId)
	}

	// carry over current logger (if not present, apmtracing.StartTransaction will build on the default logger)
	logger := zerolog.Ctx(sourceCtx)
	if logger != nil {
		ctx = logger.WithContext(ctx)
	}

	// start APM transaction, so we get a trace id, and add APM enabled logger to the context
	var apmCancel context.CancelFunc
	ctx, apmCancel = apmtracing.StartTransaction(ctx, name, traceTransactionType)

	// log a line to link parent to child transaction
	aulogging.Logger.Ctx(sourceCtx).Info().Printf("starting asynchronous transaction %s of type %s, trace.id %s",
		name, traceTransactionType, auapmlogging.ExtractTraceId(ctx))

	cancel = func() {
		apmCancel()
		fullCancel()
	}

	return ctx, cancel
}

// AsyncProcessingChildContext creates a child context for asynchronous processing that is started
// from a previously running context.
//
// The difference to AsyncCopyRequestContext is that here the original context is expected to remain active,
// Meaning, its cancellation has to wait for the subroutine to complete.
//
// We also create a tracing span.
//
// name is used as the name for the trace span.
//
// traceSpanType is used as the trace span type. You should limit the number of span types you use.
// See https://www.elastic.co/guide/en/apm/guide/current/data-model-transactions.html for details.
//
// Usage example:
//
//	ctx, cancel := contexthelper.AsyncProcessingSubContext(r.Context(), "Metadata Request Info X", "app.internal.update")
//	defer cancel()
//
// The returned cancel function will close the span, clean up, then finally cancel the context.
func AsyncProcessingChildContext(sourceCtx context.Context, name string, traceSpanType string) (ctx context.Context, cancel context.CancelFunc) {
	// here, the context IS a child
	ctx = sourceCtx

	// child context so we do not cancel global context
	var fullCancel context.CancelFunc
	ctx, fullCancel = context.WithCancel(ctx)

	// carry over Zipkin request id
	requestId := requestid.GetReqID(sourceCtx)
	if requestId != "" {
		ctx = requestid.PutReqID(ctx, requestId)
	}

	// carry over APM transaction, if exists
	tx := apm.TransactionFromContext(sourceCtx)
	if tx != nil {
		ctx = apm.ContextWithTransaction(ctx, tx)
	}

	// carry over current logger, or use global fallback (better than no logger at all)
	logger := zerolog.Ctx(sourceCtx)
	if logger != nil {
		ctx = logger.WithContext(ctx)
	} else {
		ctx = log.Logger.WithContext(ctx)
	}

	// start APM span and update logger in context
	apmCancel := func() {} // default cancel does nothing
	if tx != nil {
		ctx, apmCancel = apmtracing.StartSpan(ctx, name, traceSpanType)
	}

	cancel = func() {
		apmCancel()
		fullCancel()
	}

	return ctx, cancel
}
