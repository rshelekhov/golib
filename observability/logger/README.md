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
    cfg := logger.Config{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        Env:            "development",
        Level:          slog.LevelDebug,
        // Endpoint empty = stdout exporter
    }
    loggerProvider, otelLogger, err := logger.Init(context.Background(), cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer loggerProvider.Shutdown(context.Background())

    // Initialize tracing for correlation
    tracingCfg := tracing.Config{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        Env:            "development",
        ExporterType:   tracing.ExporterStdout,
    }
    tracerProvider, err := tracing.Init(context.Background(), tracingCfg)
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
    cfg := logger.Config{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        Env:            "production",
        Level:          slog.LevelInfo,
        Endpoint:       "localhost:4317", // OTLP endpoint
    }
    loggerProvider, otelLogger, err := logger.Init(ctx, cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer loggerProvider.Shutdown(context.Background())

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
cfg := logger.Config{
    ServiceName:    "high-traffic-service",
    ServiceVersion: "1.0.0",
    Env:            "production",
    Level:          slog.LevelWarn,
}
loggerProvider, logger, err := logger.Init(context.Background(), cfg)

// Error-only logging for critical services
cfg = logger.Config{
    ServiceName:    "critical-service",
    ServiceVersion: "1.0.0",
    Env:            "production",
    Level:          slog.LevelError,
    Endpoint:       "localhost:4317",
}
loggerProvider, logger, err = logger.Init(ctx, cfg)
```

## API Reference

### Configuration

```go
type Config struct {
    ServiceName    string
    ServiceVersion string
    Env            string
    Level          slog.Level
    Endpoint       string // OTLP endpoint. If empty, stdout exporter is used.
}
```

### Initialization Function

#### `Init(ctx context.Context, cfg Config) (*log.LoggerProvider, *slog.Logger, error)`

Initializes logger with automatic exporter selection based on configuration.

**Exporter Selection:**

- If `Endpoint` is empty: stdout exporter (development)
- If `Endpoint` is set: OTLP exporter (production)

**Parameters:**

- `ctx` - Context for initialization
- `cfg` - Logger configuration

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
cfg := logger.Config{
    ServiceName: "service", ServiceVersion: "1.0.0", Env: "local", Level: slog.LevelDebug,
}

// Production - info and above
cfg = logger.Config{
    ServiceName: "service", ServiceVersion: "1.0.0", Env: "prod", Level: slog.LevelInfo, Endpoint: "localhost:4317",
}

// High-traffic - warnings and errors only
cfg = logger.Config{
    ServiceName: "service", ServiceVersion: "1.0.0", Env: "prod", Level: slog.LevelWarn, Endpoint: "localhost:4317",
}
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

## Examples

See [examples/main.go](examples/main.go) for complete usage examples.
