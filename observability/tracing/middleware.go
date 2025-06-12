package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/stats"
)

// HTTPMiddleware returns middleware for HTTP with tracing
func HTTPMiddleware(handler http.Handler, serviceName string) http.Handler {
	return otelhttp.NewHandler(handler, serviceName)
}

// GRPCServerStatsHandler returns stats.Handler for gRPC server with tracing
func GRPCServerStatsHandler() stats.Handler {
	return otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(otel.GetTracerProvider()))
}

// GRPCClientStatsHandler returns stats.Handler for gRPC client with tracing
func GRPCClientStatsHandler() stats.Handler {
	return otelgrpc.NewClientHandler(otelgrpc.WithTracerProvider(otel.GetTracerProvider()))
}
