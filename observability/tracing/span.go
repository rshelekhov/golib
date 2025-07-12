package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Export types for convenience
// Use tracing.Attribute, tracing.SpanKind, tracing.SpanStartOption
//nolint:revive

type Attribute = attribute.KeyValue
type SpanKind = trace.SpanKind
type SpanStartOption = trace.SpanStartOption

const SpanKindClient = trace.SpanKindClient
const SpanKindServer = trace.SpanKindServer

const tracerName = "github.com/rshelekhov/golib/observability/tracing"

// StartSpan creates a span with an arbitrary name
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := otel.Tracer(tracerName)
	return tracer.Start(ctx, name, opts...)
}

// SpanFromHTTP creates a span for an HTTP request
func SpanFromHTTP(ctx context.Context, method, path string) (context.Context, trace.Span) {
	tracer := otel.Tracer(tracerName)
	return tracer.Start(ctx, method+" "+path,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("http.method", method),
			attribute.String("http.route", path),
		),
	)
}

// SpanFromGRPC creates a span for a gRPC method
func SpanFromGRPC(ctx context.Context, method string) (context.Context, trace.Span) {
	tracer := otel.Tracer(tracerName)
	return tracer.Start(ctx, method,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("rpc.service", "grpc"),
			attribute.String("rpc.method", method),
		),
	)
}

// OutgoingSpan creates a span for outgoing calls (DB, external, etc)
func OutgoingSpan(ctx context.Context, name string, spanKind SpanKind, attrs ...Attribute) (context.Context, trace.Span) {
	tracer := otel.Tracer(tracerName)
	return tracer.Start(ctx, name,
		trace.WithSpanKind(spanKind),
		trace.WithAttributes(attrs...),
	)
}

// RecordError records the provided error on the span and sets the span status to codes.Error.
// It is safe to call with nil span or error.
func RecordError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
