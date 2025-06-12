# Logger

OpenTelemetry-integrated structured logging with automatic trace correlation.

## Features

- **OpenTelemetry Logs API** - Full OTel logs integration
- **Automatic trace correlation** - Logs automatically include trace_id and span_id
- **Structured logging** - Built on Go's `log/slog`
- **Multiple exporters** - Stdout for development, OTLP for production
- **Batched processing** - Efficient log delivery
- **Resource attributes** - Service name, version, environment

## Quick Start

### Development (Stdout)

```go
package main

import (
    "context"
    "log"

    "github.com/rshelekhov/golib/observability/logger"
    "github.com/rshelekhov/golib/observability/tracing"
)

func main() {
    // Initialize logger
    loggerProvider, otelLogger, err := logger.InitLogger("my-service", "1.0.0", "development")
    if err != nil {
        log.Fatal(err)
    }
    defer loggerProvider.Shutdown(context.Background())

    // Initialize tracing for correlation
    tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
    if err != nil {
        log.Fatal(err)
    }
    defer tracerProvider.Shutdown(context.Background())

    // Create span and log within it
    tracer := tracerProvider.Tracer("example")
    ctx, span := tracer.Start(context.Background(), "operation")
    defer span.End()

    // Logs automatically include trace_id and span_id
    otelLogger.InfoContext(ctx, "processing started", "user_id", "123")
    otelLogger.ErrorContext(ctx, "processing failed", "error", "timeout")
}
```

### Production (OTLP)

```go
func main() {
    ctx := context.Background()

    // Initialize logger with OTLP
    loggerProvider, otelLogger, err := logger.InitLoggerOTLP(
        ctx, "my-service", "1.0.0", "production", "localhost:4317")
    if err != nil {
        log.Fatal(err)
    }
    defer loggerProvider.Shutdown(context.Background())

    // Initialize tracing with OTLP
    tracerProvider, err := tracing.InitTracerOTLP(
        ctx, "my-service", "1.0.0", "production", "localhost:4317")
    if err != nil {
        log.Fatal(err)
    }
    defer tracerProvider.Shutdown(context.Background())

    // Use logger with automatic correlation
    tracer := tracerProvider.Tracer("example")
    ctx, span := tracer.Start(ctx, "operation")
    defer span.End()

    otelLogger.InfoContext(ctx, "OTLP processing", "user_id", "456")
}
```

## API Reference

### Initialization Functions

#### `InitLogger(serviceName, serviceVersion, env string) (*log.LoggerProvider, *slog.Logger, error)`

Initializes logger with stdout exporter for development.

**Parameters:**

- `serviceName` - Service name for resource attributes
- `serviceVersion` - Service version for resource attributes
- `env` - Environment (development, staging, production)

**Returns:**

- `*log.LoggerProvider` - For shutdown management
- `*slog.Logger` - OTel-integrated slog logger
- `error` - Initialization error

#### `InitLoggerOTLP(ctx context.Context, serviceName, serviceVersion, env, endpoint string) (*log.LoggerProvider, *slog.Logger, error)`

Initializes logger with OTLP exporter for production.

**Parameters:**

- `ctx` - Context for initialization
- `serviceName` - Service name for resource attributes
- `serviceVersion` - Service version for resource attributes
- `env` - Environment (development, staging, production)
- `endpoint` - OTLP collector endpoint (e.g., "localhost:4317")

**Returns:**

- `*log.LoggerProvider` - For shutdown management
- `*slog.Logger` - OTel-integrated slog logger
- `error` - Initialization error

### Logger Usage

The returned `*slog.Logger` is a standard Go slog logger with OTel integration:

```go
// Standard slog methods with automatic trace correlation
logger.InfoContext(ctx, "message", "key", "value")
logger.WarnContext(ctx, "warning", "field", "value")
logger.ErrorContext(ctx, "error", "error", err)
logger.DebugContext(ctx, "debug info", "data", data)
```

## Key Benefits

### 1. Automatic Trace Correlation

When logging within an active span, logs automatically include:

- `trace_id` - Links log to distributed trace
- `span_id` - Links log to specific span
- Service resource attributes

### 2. Unified Observability

All three signals (logs, traces, metrics) use the same:

- OTLP endpoint
- Resource attributes
- Correlation context

### 3. Production Ready

- **Batched processing** - Logs sent in efficient batches
- **Async delivery** - Non-blocking log processing
- **Resource attributes** - Proper service identification
- **Error handling** - Graceful degradation

### 4. Standard slog Interface

No custom logging methods - uses standard Go `log/slog` API with OTel enhancement.

## Architecture

```
Application Code
       ↓
   slog.Logger (OTel bridge)
       ↓
   OTel LoggerProvider
       ↓
   BatchProcessor
       ↓
   Exporter (Stdout/OTLP)
       ↓
   Backend (Console/Collector)
```

## Integration with Other Components

### With Tracing

```go
// Initialize both
tracerProvider, _ := tracing.InitTracer("service", "1.0.0", "dev")
loggerProvider, logger, _ := logger.InitLogger("service", "1.0.0", "dev")

// Use together
tracer := tracerProvider.Tracer("example")
ctx, span := tracer.Start(context.Background(), "operation")
defer span.End()

// Log is automatically correlated with span
logger.InfoContext(ctx, "operation completed")
```

### With Observability Facade

```go
obs, err := observability.Init(context.Background(), observability.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:           "development",
    EnableMetrics: true,
})

// obs.Logger is ready to use with trace correlation
obs.Logger.InfoContext(ctx, "request processed")
```

## Examples

See `examples/main.go` for complete working examples of both stdout and OTLP configurations.
