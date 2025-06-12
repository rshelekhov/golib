package metrics

import (
	"context"
	"time"
	"sync"
	"log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	grpcRequestsCounter metric.Int64Counter
	grpcLatencyHistogram metric.Float64Histogram
	initGRPCMetricsOnce sync.Once
)

func initGRPCMetrics() {
	initGRPCMetricsOnce.Do(func() {
		meter := OtelMeter()
		var err error
		
		grpcRequestsCounter, err = meter.Int64Counter(
			"grpc_server_requests_total",
			metric.WithDescription("Total number of gRPC requests received."),
		)
		if err != nil {
			log.Fatalf("failed to create grpc_server_requests_total counter: %v", err)
		}
		
		grpcLatencyHistogram, err = meter.Float64Histogram(
			"grpc_server_handling_seconds",
			metric.WithDescription("gRPC request handling duration in seconds."),
		)
		if err != nil {
			log.Fatalf("failed to create grpc_server_handling_seconds histogram: %v", err)
		}
	})
}

// UnaryServerInterceptor returns grpc.UnaryServerInterceptor for otel metrics
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	initGRPCMetrics()
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err).String()
		service, method := splitMethod(info.FullMethod)
		grpcRequestsCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("service", service),
			attribute.String("method", method),
			attribute.String("code", code),
		))
		grpcLatencyHistogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String("service", service),
			attribute.String("method", method),
		))
		return resp, err
	}
}

// StreamServerInterceptor returns grpc.StreamServerInterceptor for otel metrics
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	initGRPCMetrics()
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		err := handler(srv, ss)
		code := status.Code(err).String()
		service, method := splitMethod(info.FullMethod)
		grpcRequestsCounter.Add(ss.Context(), 1, metric.WithAttributes(
			attribute.String("service", service),
			attribute.String("method", method),
			attribute.String("code", code),
		))
		grpcLatencyHistogram.Record(ss.Context(), time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String("service", service),
			attribute.String("method", method),
		))
		return err
	}
}

func splitMethod(fullMethod string) (service, method string) {
	// fullMethod: /package.service/method
	if len(fullMethod) == 0 || fullMethod[0] != '/' {
		return "", ""
	}
	fullMethod = fullMethod[1:]
	parts := make([]string, 2)
	for i, s := range []rune(fullMethod) {
		if s == '/' {
			parts[0] = string([]rune(fullMethod)[:i])
			parts[1] = string([]rune(fullMethod)[i+1:])
			return parts[0], parts[1]
		}
	}
	return fullMethod, ""
}
