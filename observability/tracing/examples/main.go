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
	// Standard pattern from observability courses
	tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()
	
	fmt.Println("Tracer initialized with stdout exporter")
	// ... your application logic
}

// Example of tracing initialization (OTLP)
func ExampleInitTracerOTLP() {
	ctx := context.Background()
	
	// OTLP exporter for production
	tracerProvider, err := tracing.InitTracerOTLP(ctx, "my-service", "1.0.0", "production", "localhost:4317")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()
	
	fmt.Println("Tracer initialized with OTLP exporter")
	// ... your application logic
}

// Example of HTTP middleware
func ExampleHTTPMiddleware() {
	// Initialize tracer first
	tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
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
	// Initialize tracer first
	tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
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
	// ExampleHTTPMiddleware()
	// ExampleGRPCStatsHandler()
	
	// Or run simple examples:
	ExampleSpanFromHTTP()
	ExampleSpanFromGRPC()
	ExampleOutgoingSpan()
}
