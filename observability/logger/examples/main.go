package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/rshelekhov/golib/observability/logger"
	"github.com/rshelekhov/golib/observability/tracing"
)

func main() {
	// Example 1: Stdout logger
	stdoutExample()

	// Example 2: OTLP logger
	// otlpExample()

	// Example 3: Pretty logger for local development
	prettyExample()
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

func prettyExample() {
	ctx := context.Background()

	// Initialize logger with pretty handler for local development
	loggerCfg := logger.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "local",
		Level:          slog.LevelDebug,
	}

	// Note: LoggerProvider will be nil for local env since we use pretty handler
	_, prettyLogger, err := logger.Init(ctx, loggerCfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== Pretty Logger Example ===")

	// Test different log levels with attributes
	prettyLogger.Debug("Debug message", "user_id", 123, "action", "login")
	prettyLogger.Info("User logged in", "user_id", 123, "email", "user@example.com", "ip", "192.168.1.1")
	prettyLogger.Warn("Rate limit approaching", "user_id", 123, "requests", 95, "limit", 100)
	prettyLogger.Error("Database connection failed", "error", "connection timeout", "retries", 3)

	// Example with structured data
	prettyLogger.Info("Order processed",
		"order_id", "ord_123456",
		"customer_id", 789,
		"amount", 99.99,
		"currency", "USD",
		"items", []string{"item1", "item2", "item3"},
		"timestamp", time.Now(),
	)

	// Example with grouped attributes
	groupedLogger := prettyLogger.With("service", "payment", "version", "v2.1.0")
	groupedLogger.Info("Payment processed", "transaction_id", "tx_789", "status", "success")
	groupedLogger.Error("Payment failed", "transaction_id", "tx_790", "reason", "insufficient funds")

	fmt.Println("=== End Pretty Logger Example ===")
}
