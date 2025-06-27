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
	// Initialize logger with stdout exporter
	loggerProvider, otelLogger, err := logger.InitLoggerStdout("my-service", "1.0.0", "development", slog.LevelDebug)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down logger: %v", err)
		}
	}()

	// Initialize tracing to see correlation
	tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
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
	loggerProvider, otelLogger, err := logger.InitLoggerOTLP(ctx, "my-service", "1.0.0", "production", "localhost:4317", slog.LevelInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down logger: %v", err)
		}
	}()

	// Initialize tracing with OTLP
	tracerProvider, err := tracing.InitTracerOTLP(ctx, "my-service", "1.0.0", "production", "localhost:4317")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
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
