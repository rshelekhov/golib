package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/rshelekhov/golib/observability"
	"github.com/rshelekhov/golib/observability/logger"
	"github.com/rshelekhov/golib/observability/metrics"
	"github.com/rshelekhov/golib/observability/tracing"
	"google.golang.org/grpc"
)

func main() {
	// Example 1: Local development with default debug logging
	localExample()

	// Example 2: Production setup with default info logging
	// prodExample()

	// Example 3: Custom log levels - override environment defaults
	// customLogLevelExample()

	// Example 4: Error handling for invalid configurations
	// errorHandlingExample()

	// Example 5: Manual initialization
	// manualExample()
}

// Example 1: Local development with default debug logging
func localExample() {
	// Using simplified API - debug level by default for local, metrics always disabled
	cfg, err := observability.NewConfig(observability.ConfigParams{
		Env:            observability.EnvLocal,
		ServiceName:    "my-service",
		ServiceVersion: "1.2.3",
		EnableMetrics:  false,
		OTLPEndpoint:   "",
	})
	if err != nil {
		log.Fatal(err)
	}

	obs, err := observability.Init(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := obs.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	// Simple HTTP server (no metrics - they're disabled for local)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Debug level logs will be shown in local development
		obs.Logger.DebugContext(r.Context(), "processing request", "path", r.URL.Path)
		obs.Logger.InfoContext(r.Context(), "handling request")

		if _, err := w.Write([]byte("Hello from local development!")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}))

	// gRPC server (no metrics interceptor)
	go func() {
		server := grpc.NewServer(
			grpc.StatsHandler(tracing.GRPCServerStatsHandler()),
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

	log.Printf("Traces printed to stdout with DEBUG level")
	log.Printf("Metrics are disabled for local development")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Example 2: Production setup with default info logging
//
//nolint:unused // Example function
func prodExample() {
	// Production with default info level
	cfg, err := observability.NewConfig(observability.ConfigParams{
		Env:               observability.EnvProd,
		ServiceName:       "my-service",
		ServiceVersion:    "1.2.3",
		EnableMetrics:     true,
		OTLPEndpoint:      "localhost:4317",
		OTLPTransportType: "grpc",
	})
	if err != nil {
		log.Fatal(err)
	}

	obs, err := observability.Init(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := obs.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	// No metrics HTTP endpoint needed - push model with OTLP
	http.Handle("/", metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Debug logs won't be shown in production (info level)
		obs.Logger.DebugContext(r.Context(), "this won't be logged")
		obs.Logger.InfoContext(r.Context(), "handling request")

		// Custom business metric
		metrics.IncBusinessError("validation", "invalid_input")

		if _, err := w.Write([]byte("Hello from production!")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})))

	log.Printf("Traces and metrics pushed to OTLP collector at localhost:4317")
	log.Printf("Using INFO level logging")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Example 3: Custom log levels - override environment defaults
//
//nolint:unused // Example function
func customLogLevelExample() {
	// Override local environment to use warn level instead of debug
	cfg, err := observability.NewConfig(
		observability.ConfigParams{
			Env:            observability.EnvLocal,
			ServiceName:    "my-service",
			ServiceVersion: "1.0.0",
			EnableMetrics:  false,
			OTLPEndpoint:   "",
		},
		observability.WithLogLevel(slog.LevelWarn), // Override default debug level
	)
	if err != nil {
		log.Fatal(err)
	}

	obs, err := observability.Init(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := obs.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	// No metrics endpoint - disabled for local development
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// These won't be logged (below warn level)
		obs.Logger.DebugContext(r.Context(), "debug message")
		obs.Logger.InfoContext(r.Context(), "info message")

		// This will be logged
		obs.Logger.WarnContext(r.Context(), "warning message", "path", r.URL.Path)

		if _, err := w.Write([]byte("Custom log level example!")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}))

	log.Printf("Using WARN level - only warnings and errors will be logged")
	log.Printf("Metrics are disabled for local development")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Example 4: Error handling for invalid configurations
//
//nolint:unused // Example function
func errorHandlingExample() {
	// This will fail - unknown environment
	cfg, err := observability.NewConfig(observability.ConfigParams{
		Env:               "staging",
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "localhost:4317",
		OTLPTransportType: "grpc",
	})
	if err != nil {
		log.Printf("Expected error: %v", err)
		// Output: unsupported environment: staging (supported: local, dev, prod)
	}

	// This will fail - missing OTLP endpoint for prod
	cfg, err = observability.NewConfig(observability.ConfigParams{
		Env:            observability.EnvProd,
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		EnableMetrics:  true,
		OTLPEndpoint:   "", // Missing endpoint
	})
	if err != nil {
		log.Printf("Expected error: %v", err)
		// Output: OTLP endpoint is required for environment prod
	}

	// This will succeed
	cfg, err = observability.NewConfig(observability.ConfigParams{
		Env:               observability.EnvProd,
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "localhost:4317",
		OTLPTransportType: "grpc",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Valid config created: %+v", cfg)
}

// Example 5: Manual initialization for fine control
//
//nolint:unused // Example function
func manualExample() {
	ctx := context.Background()

	// Initialize tracing manually
	tracingCfg := tracing.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "development",
		ExporterType:   tracing.ExporterStdout,
	}
	tracerProvider, err := tracing.Init(ctx, tracingCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	// Initialize logger manually
	loggerCfg := logger.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "development",
		Level:          slog.LevelDebug,
		Endpoint:       "", // Use stdout
	}
	loggerProvider, otelLogger, err := logger.Init(ctx, loggerCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := loggerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down logger: %v", err)
		}
	}()

	// Initialize metrics manually
	metricsCfg := metrics.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "development",
		ExporterType:   metrics.ExporterPrometheus,
	}
	meterProvider, handler, err := metrics.Init(ctx, metricsCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down metrics: %v", err)
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
		// Use the custom logger
		otelLogger.InfoContext(r.Context(), "handling request", "path", r.URL.Path)

		// Increment custom counter
		counter.Add(r.Context(), 1)

		if _, err := w.Write([]byte("Manual setup complete!")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})))

	log.Printf("Manual setup with all components initialized separately")
	log.Printf("Metrics available at http://localhost:8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
