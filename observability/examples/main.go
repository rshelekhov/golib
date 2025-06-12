package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/rshelekhov/golib/observability"
	"github.com/rshelekhov/golib/observability/logger"
	"github.com/rshelekhov/golib/observability/metrics"
	"github.com/rshelekhov/golib/observability/tracing"
	"google.golang.org/grpc"
)

func main() {
	// Example 1: Stdout tracing + Prometheus metrics
	stdoutExample()
	
	// Example 2: OTLP tracing + OTLP metrics
	// otlpExample()
	
	// Example 3: Manual initialization
	// manualExample()
}

// Example 1: Simple initialization with stdout tracing and Prometheus metrics
func stdoutExample() {
	obs, err := observability.Init(context.Background(), observability.Config{
		Env:            "development",
		ServiceName:    "my-service",
		ServiceVersion: "1.2.3",
		EnableMetrics:  true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := obs.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	// Setup HTTP server with metrics endpoint
	http.Handle("/metrics", obs.MetricsHandler) // Prometheus scrapes this
	http.Handle("/", metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obs.Logger.InfoContext(r.Context(), "handling request")
		if _, err := w.Write([]byte("Hello from stdout setup!")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})))

	// gRPC server
	go func() {
		server := grpc.NewServer(
			grpc.StatsHandler(tracing.GRPCServerStatsHandler()),
			grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()),
		)

		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		log.Printf("gRPC server listening on :50051")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Printf("Traces printed to stdout")
	log.Printf("Prometheus metrics available at http://localhost:8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Example 2: OTLP exporters for production
//nolint:unused // Example function
func otlpExample() {
	obs, err := observability.InitWithOTLP(context.Background(), observability.Config{
		Env:            "production",
		ServiceName:    "my-service",
		ServiceVersion: "1.2.3",
		EnableMetrics:  true,
	}, "localhost:4317") // OTLP endpoint for both tracing and metrics
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := obs.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	// No metrics HTTP endpoint needed - push model
	http.Handle("/", metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obs.Logger.InfoContext(r.Context(), "handling request")
		
		// Custom business metric
		metrics.IncBusinessError("validation", "invalid_input")
		
		if _, err := w.Write([]byte("Hello from OTLP setup!")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})))

	log.Printf("Traces and metrics pushed to OTLP collector at localhost:4317")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Example 3: Manual initialization for fine control
//nolint:unused // Example function
func manualExample() {
	// Initialize tracing manually
	tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	// Initialize metrics manually
	meterProvider, handler, err := metrics.InitMeter("my-service", "1.0.0", "development")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down meter: %v", err)
		}
	}()

	// Setup logger manually
	loggerProvider, logger, err := logger.InitLogger("my-service", "1.0.0", "development")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down logger: %v", err)
		}
	}()

	// Create custom metrics
	meter := metrics.OtelMeter()
	counter, err := meter.Int64Counter("manual_counter")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", handler)
	http.Handle("/", metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "handling request manually")
		
		// Increment custom counter
		counter.Add(r.Context(), 1)
		
		if _, err := w.Write([]byte("Hello from manual setup!")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})))

	log.Printf("Manual setup complete")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
