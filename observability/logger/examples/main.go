package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/rshelekhov/golib/observability/logger"
	"github.com/rshelekhov/golib/observability/tracing"
)

func main() {
	// Example 1: Stdout logger
	stdoutExample()

	// Example 2: OTLP logger
	// otlpExample()
}

func stdoutExample() {
	ctx := context.Background()

	// Initialize logger with stdout exporter
	loggerCfg := logger.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "development",
		Level:          slog.LevelDebug,
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

	// Initialize tracing to see correlation
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

	// Create a span and log within it
	tracer := tracerProvider.Tracer("example")
	ctx, span := tracer.Start(context.Background(), "example_operation")
	defer span.End()

	// These logs will automatically include trace_id and span_id
	otelLogger.InfoContext(ctx, "processing started", "user_id", "123")
	otelLogger.WarnContext(ctx, "validation warning", "field", "email")
	otelLogger.ErrorContext(ctx, "processing failed", "error", "database timeout")

	fmt.Println("Check stdout for correlated logs and traces!")
}

//nolint:unused // Example function
func otlpExample() {
	ctx := context.Background()

	// Initialize logger with OTLP exporter
	loggerCfg := logger.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "production",
		Level:          slog.LevelInfo,
		Endpoint:       "localhost:4317",
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

	// Initialize tracing with OTLP
	tracingCfg := tracing.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "production",
		ExporterType:   tracing.ExporterOTLP,
		OTLPEndpoint:   "localhost:4317",
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

	// Create a span and log within it
	tracer := tracerProvider.Tracer("example")
	ctx, span := tracer.Start(ctx, "otlp_operation")
	defer span.End()

	// These logs will be sent to OTLP collector with trace correlation
	otelLogger.InfoContext(ctx, "OTLP processing started", "user_id", "456")
	otelLogger.ErrorContext(ctx, "OTLP processing failed", "error", "network timeout")

	fmt.Println("Logs and traces sent to OTLP collector!")
}
