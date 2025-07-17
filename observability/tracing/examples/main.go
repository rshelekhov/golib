package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/rshelekhov/golib/observability/tracing"
	"google.golang.org/grpc"
)

// Example of tracing initialization (stdout)
func ExampleInitTracer() {
	ctx := context.Background()

	// Standard pattern using new Config API
	cfg := tracing.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "development",
		ExporterType:   tracing.ExporterStdout,
	}
	tracerProvider, err := tracing.Init(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	fmt.Println("Tracer initialized with stdout exporter")
	// ... your application logic
}

// Example of tracing initialization (OTLP with TLS)
func ExampleInitTracerOTLP() {
	ctx := context.Background()

	// OTLP exporter for production (with TLS by default)
	cfg := tracing.Config{
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		Env:               "production",
		ExporterType:      tracing.ExporterOTLP,
		OTLPEndpoint:      "otel-collector.company.com:4317",
		OTLPTransportType: tracing.OTLPGRPC,
		OTLPInsecure:      false, // Uses TLS (default for production)
	}
	tracerProvider, err := tracing.Init(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	fmt.Println("Tracer initialized with OTLP exporter using TLS")
	// ... your application logic
}

// Example of tracing initialization (OTLP without TLS for development)
func ExampleInitTracerOTLPInsecure() {
	ctx := context.Background()

	// OTLP exporter for development (insecure connection)
	cfg := tracing.Config{
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		Env:               "development",
		ExporterType:      tracing.ExporterOTLP,
		OTLPEndpoint:      "localhost:4317",
		OTLPTransportType: tracing.OTLPGRPC,
		OTLPInsecure:      true, // Uses insecure connection (default for dev)
	}
	tracerProvider, err := tracing.Init(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	fmt.Println("Tracer initialized with OTLP exporter without TLS (development)")
	// ... your application logic
}

// Example of tracing initialization (OTLP HTTP transport with TLS)
func ExampleInitTracerOTLPHTTP() {
	ctx := context.Background()

	// OTLP HTTP exporter for production (with TLS)
	cfg := tracing.Config{
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		Env:               "production",
		ExporterType:      tracing.ExporterOTLP,
		OTLPEndpoint:      "https://otel-collector.company.com:4318",
		OTLPTransportType: tracing.OTLPHTTP,
		OTLPInsecure:      false, // Uses TLS (default for production)
	}
	tracerProvider, err := tracing.Init(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	fmt.Println("Tracer initialized with OTLP HTTP exporter using TLS")
	// ... your application logic
}

// Example of HTTP middleware
func ExampleHTTPMiddleware() {
	ctx := context.Background()

	// Initialize tracer first
	cfg := tracing.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "development",
		ExporterType:   tracing.ExporterStdout,
	}
	tracerProvider, err := tracing.Init(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	handler := tracing.HTTPMiddleware(mux, "my-service")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// Example of gRPC server with tracing
func ExampleGRPCStatsHandler() {
	ctx := context.Background()

	// Initialize tracer first
	cfg := tracing.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "development",
		ExporterType:   tracing.ExporterStdout,
	}
	tracerProvider, err := tracing.Init(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	server := grpc.NewServer(
		grpc.StatsHandler(tracing.GRPCServerStatsHandler()),
	)

	// Register your gRPC services here
	// pb.RegisterYourServiceServer(server, &yourServiceImpl{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on :50051")
	log.Fatal(server.Serve(lis))
}

// Example of creating span in HTTP handler
func ExampleSpanFromHTTP() {
	ctx := context.Background()
	_, span := tracing.SpanFromHTTP(ctx, "GET", "/api/v1/users/{id}")
	defer span.End()

	// Your business logic here
	fmt.Printf("Trace ID: %s\n", span.SpanContext().TraceID())
}

// Example of creating span in gRPC handler
func ExampleSpanFromGRPC() {
	ctx := context.Background()
	_, span := tracing.SpanFromGRPC(ctx, "UserService.GetUser")
	defer span.End()

	// Your business logic here
	fmt.Printf("Trace ID: %s\n", span.SpanContext().TraceID())
}

// Example of outgoing span (DB, external API)
func ExampleOutgoingSpan() {
	ctx := context.Background()
	_, span := tracing.OutgoingSpan(ctx, "db.query", tracing.SpanKindClient,
		tracing.String("db.system", "postgresql"),
		tracing.String("db.statement", "SELECT * FROM users WHERE id = ?"),
	)
	defer span.End()

	// Execute DB query
	fmt.Printf("Span ID: %s\n", span.SpanContext().SpanID())
}

func main() {
	// Run one of the examples
	// Uncomment the needed example:

	// ExampleInitTracer()
	// ExampleInitTracerOTLP()
	// ExampleInitTracerOTLPInsecure()
	// ExampleInitTracerOTLPHTTP()
	// ExampleHTTPMiddleware()
	// ExampleGRPCStatsHandler()

	// Or run simple examples:
	ExampleSpanFromHTTP()
	ExampleSpanFromGRPC()
	ExampleOutgoingSpan()
}
