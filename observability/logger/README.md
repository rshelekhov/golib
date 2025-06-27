# Logger

OpenTelemetry-integrated structured logging with automatic trace correlation and configurable log levels.

## Features

- **OpenTelemetry Logs API** - Full OTel logs integration
- **Automatic trace correlation** - Logs automatically include trace_id and span_id
- **Structured logging** - Built on Go's `log/slog`
- **Multiple exporters** - Stdout for development, OTLP for production
- **Configurable log levels** - Debug, Info, Warn, Error filtering
- **Batched processing** - Efficient log delivery
- **Resource attributes** - Service name, version, environment

## Quick Start

### Development (Stdout)

```go
package main

import (
    "context"
    "log"
    "log/slog"

    "github.com/rshelekhov/golib/observability/logger"
    "github.com/rshelekhov/golib/observability/tracing"
)

func main() {
    // Initialize logger with debug level for development
    loggerProvider, otelLogger, err := logger.InitLoggerStdout("my-service", "1.0.0", "development", slog.LevelDebug)
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

    // Debug logs will be shown
    otelLogger.DebugContext(ctx, "debug information", "user_id", "123")
    // Logs automatically include trace_id and span_id
    otelLogger.InfoContext(ctx, "processing started", "user_id", "123")
    otelLogger.ErrorContext(ctx, "processing failed", "error", "timeout")
}
```

### Production (OTLP)

```go
func main() {
    ctx := context.Background()

    // Initialize logger with OTLP and info level for production
    loggerProvider, otelLogger, err := logger.InitLoggerOTLP(
        ctx, "my-service", "1.0.0", "production", "localhost:4317", slog.LevelInfo)
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

    // Debug logs won't be shown (info level)
    otelLogger.DebugContext(ctx, "this won't be logged")
    otelLogger.InfoContext(ctx, "OTLP processing", "user_id", "456")
}
```

### Custom Log Levels

```go
// High-traffic service - only warnings and errors
loggerProvider, logger, err := logger.InitLoggerStdout("high-traffic-service", "1.0.0", "production", slog.LevelWarn)

// Error-only logging for critical services
loggerProvider, logger, err := logger.InitLoggerOTLP(ctx, "critical-service", "1.0.0", "production", "localhost:4317", slog.LevelError)
```

## API Reference

### Initialization Functions

#### `InitLoggerStdout(serviceName, serviceVersion, env string, level slog.Level) (*log.LoggerProvider, *slog.Logger, error)`

Initializes logger with stdout exporter for development.

**Parameters:**

- `serviceName` - Service name for resource attributes
- `serviceVersion` - Service version for resource attributes
- `env` - Environment (development, staging, production)
- `level` - Minimum log level (slog.LevelDebug, LevelInfo, LevelWarn, LevelError)

**Returns:**

- `*log.LoggerProvider` - For shutdown management
- `*slog.Logger` - OTel-integrated slog logger
- `error` - Initialization error

#### `InitLoggerOTLP(ctx context.Context, serviceName, serviceVersion, env, endpoint string, level slog.Level) (*log.LoggerProvider, *slog.Logger, error)`

Initializes logger with OTLP exporter for production.

**Parameters:**

- `ctx` - Context for initialization
- `serviceName` - Service name for resource attributes
- `serviceVersion` - Service version for resource attributes
- `env` - Environment (development, staging, production)
- `endpoint` - OTLP collector endpoint (e.g., "localhost:4317")
- `level` - Minimum log level (slog.LevelDebug, LevelInfo, LevelWarn, LevelError)

**Returns:**

- `*log.LoggerProvider` - For shutdown management
- `*slog.Logger` - OTel-integrated slog logger
- `error` - Initialization error

### Log Levels

- `slog.LevelDebug` (-4) - All logs, use for local development
- `slog.LevelInfo` (0) - Info and above, production default
- `slog.LevelWarn` (4) - Warnings and errors only, high-traffic services
- `slog.LevelError` (8) - Errors only, critical services

### Logger Usage

The returned `*slog.Logger` is a standard Go slog logger with OTel integration:

```go
// Standard slog methods with automatic trace correlation
logger.DebugContext(ctx, "debug info", "data", data)        // Level -4
logger.InfoContext(ctx, "message", "key", "value")          // Level 0
logger.WarnContext(ctx, "warning", "field", "value")        // Level 4
logger.ErrorContext(ctx, "error", "error", err)             // Level 8
```

## Key Benefits

### 1. Automatic Trace Correlation

When logging within an active span, logs automatically include:

- `trace_id` - Links log to distributed trace
- `span_id` - Links log to specific span
- Service resource attributes

### 2. Configurable Log Levels

Different log levels for different environments:

```go
// Local development - see everything
InitLoggerStdout("service", "1.0.0", "local", slog.LevelDebug)

// Production - info and above
InitLoggerOTLP(ctx, "service", "1.0.0", "prod", endpoint, slog.LevelInfo)

// High-traffic - warnings and errors only
InitLoggerOTLP(ctx, "service", "1.0.0", "prod", endpoint, slog.LevelWarn)
```

### 3. Unified Observability

All three signals (logs, traces, metrics) use the same:

- OTLP endpoint
- Resource attributes
- Correlation context

### 4. Production Ready

- **Batched processing** - Logs sent in efficient batches
- **Async delivery** - Non-blocking log processing
- **Resource attributes** - Proper service identification
- **Level filtering** - Efficient log processing
- **Error handling** - Graceful degradation

### 5. Standard slog Interface

No custom logging methods - uses standard Go `log/slog` API with OTel enhancement.

## Architecture

```
Application Code
       ↓
   slog.Logger (OTel bridge + level filter)
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
loggerProvider, logger, _ := logger.InitLoggerStdout("service", "1.0.0", "dev", slog.LevelDebug)

// Use together
tracer := tracerProvider.Tracer("example")
ctx, span := tracer.Start(context.Background(), "operation")
defer span.End()

// Log is automatically correlated with span
logger.InfoContext(ctx, "operation completed")
```

### With Observability Facade

```go
// Using helper functions
cfg := observability.NewLocalConfig("my-service", "1.0.0", true)  // debug level
cfg := observability.NewProdConfig("my-service", "1.0.0", "localhost:4317", true)  // info level

obs, err := observability.Setup(context.Background(), cfg)

// obs.Logger is ready to use with trace correlation
obs.Logger.InfoContext(ctx, "request processed")
```

## Examples

### Different Log Levels in Action

```go
logger, _, _ := logger.InitLoggerStdout("test-service", "1.0.0", "dev", slog.LevelWarn)

// These won't be logged (below warn level)
logger.DebugContext(ctx, "debug message")
logger.InfoContext(ctx, "info message")

// These will be logged
logger.WarnContext(ctx, "warning message")
logger.ErrorContext(ctx, "error message")
```

### Kubernetes Deployment

```go
// For Kubernetes, use stdout with appropriate level
logger, _, _ := logger.InitLoggerStdout("k8s-service", "1.0.0", "production", slog.LevelInfo)

// Logs go to stdout and are collected by k8s logging
logger.InfoContext(ctx, "service started")
```

See `examples/main.go` for complete working examples of different configurations and log levels.
