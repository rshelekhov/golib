# Logger

OpenTelemetry-integrated structured logging with automatic trace correlation and configurable log levels.

## Features

- **OpenTelemetry Logs API** - Full OTel logs integration
- **Pretty local logging** - Colorized, human-readable output for local development
- **Automatic trace correlation** - Logs automatically include trace_id and span_id
- **Structured logging** - Built on Go's `log/slog`
- **Multiple exporters** - Pretty handler for local, stdout for development, OTLP for production
- **Configurable log levels** - Debug, Info, Warn, Error filtering
- **Batched processing** - Efficient log delivery
- **Resource attributes** - Service name, version, environment

## Quick Start

### Local Development (Pretty Logs)

```go
package main

import (
    "context"
    "log"
    "log/slog"

    "github.com/rshelekhov/golib/observability/logger"
)

func main() {
    // Initialize logger with pretty handler for local development
    cfg := logger.Config{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        Env:            "local",
        Level:          slog.LevelDebug,
        // Endpoint not needed for local - uses pretty handler
    }

    // Note: LoggerProvider will be nil for local env
    _, prettyLogger, err := logger.Init(context.Background(), cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Colorized output with structured data
    prettyLogger.Debug("Debug message", "user_id", 123, "action", "login")
    prettyLogger.Info("User logged in", "user_id", 123, "email", "user@example.com")
    prettyLogger.Warn("Rate limit approaching", "user_id", 123, "requests", 95)
    prettyLogger.Error("Database connection failed", "error", "connection timeout")
}
```

Output:

```
[13:45:05.123] DEBUG: Debug message {
  "action": "login",
  "user_id": 123
}
[13:45:05.124] INFO: User logged in {
  "email": "user@example.com",
  "user_id": 123
}
[13:45:05.125] WARN: Rate limit approaching {
  "requests": 95,
  "user_id": 123
}
[13:45:05.126] ERROR: Database connection failed {
  "error": "connection timeout"
}
```

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

### Production (OTLP with TLS)

```go
func main() {
    ctx := context.Background()

    // Initialize logger with OTLP and info level for production (TLS by default)
    cfg := logger.Config{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        Env:            "production",
        Level:          slog.LevelInfo,
        Endpoint:       "otel-collector.company.com:4317", // OTLP endpoint
        OTLPInsecure:   false, // Uses TLS (recommended for production)
    }
    loggerProvider, otelLogger, err := logger.Init(ctx, cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer loggerProvider.Shutdown(context.Background())

    // Initialize tracing with OTLP
    tracingCfg := tracing.Config{
        ServiceName:       "my-service",
        ServiceVersion:    "1.0.0",
        Env:               "production",
        ExporterType:      tracing.ExporterOTLP,
        OTLPEndpoint:      "otel-collector.company.com:4317",
        OTLPTransportType: tracing.OTLPGRPC,
        OTLPInsecure:      false, // Uses TLS (recommended for production)
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
    otelLogger.InfoContext(ctx, "OTLP processing with TLS", "user_id", "456")
}
```

### Development (OTLP without TLS)

```go
func main() {
    ctx := context.Background()

    // Initialize logger with OTLP for development (insecure connection)
    cfg := logger.Config{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        Env:            "development",
        Level:          slog.LevelDebug,
        Endpoint:       "localhost:4317", // Local OTLP collector
        OTLPInsecure:   true, // No TLS for local development
    }
    loggerProvider, otelLogger, err := logger.Init(ctx, cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer loggerProvider.Shutdown(context.Background())

    // Initialize tracing with OTLP (insecure)
    tracingCfg := tracing.Config{
        ServiceName:       "my-service",
        ServiceVersion:    "1.0.0",
        Env:               "development",
        ExporterType:      tracing.ExporterOTLP,
        OTLPEndpoint:      "localhost:4317",
        OTLPTransportType: tracing.OTLPGRPC,
        OTLPInsecure:      true, // No TLS for local development
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

    // Debug logs will be shown (debug level)
    otelLogger.DebugContext(ctx, "debug information", "user_id", "789")
    otelLogger.InfoContext(ctx, "local OTLP processing without TLS", "user_id", "789")
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
    OTLPInsecure   bool   // If true, uses insecure OTLP connection
}
```

### Initialization Function

#### `Init(ctx context.Context, cfg Config) (*log.LoggerProvider, *slog.Logger, error)`

Initializes logger with automatic exporter selection based on configuration.

**Exporter Selection:**

- If `Env` is "local": pretty handler (colorized, human-readable)
- If `Endpoint` is empty: stdout exporter (development)
- If `Endpoint` is set: OTLP exporter (production)

**Parameters:**

- `ctx` - Context for initialization
- `cfg` - Logger configuration

**Returns:**

- `*log.LoggerProvider` - For shutdown management (nil for local env)
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

## TLS Configuration

The logger package supports configurable TLS for OTLP connections:

### TLS Configuration Examples

```go
// Production with TLS (recommended)
cfg := logger.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:            "production",
    Level:          slog.LevelInfo,
    Endpoint:       "otel-collector.company.com:4317",
    OTLPInsecure:   false, // Uses TLS (recommended for production)
}

// Development with insecure connection (local OTLP collector)
cfg := logger.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:            "development",
    Level:          slog.LevelDebug,
    Endpoint:       "localhost:4317",
    OTLPInsecure:   true, // No TLS for local development
}

// Local development (TLS not applicable - uses pretty handler)
cfg := logger.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:            "local",
    Level:          slog.LevelDebug,
    // Endpoint and OTLPInsecure not used for local env
}
```

### TLS Best Practices

- **Production**: Always use TLS (`OTLPInsecure: false`) for secure log transmission
- **Development**: Use insecure connections (`OTLPInsecure: true`) for local OTLP collectors
- **Local**: TLS configuration is not applicable (uses pretty handler, no network calls)

## Key Benefits

### 1. Pretty Local Development

For `Env: "local"`:

- **Colorized output** - Different colors for each log level
- **Human-readable format** - Easy to scan timestamps and messages
- **Structured data** - JSON-formatted attributes with proper indentation
- **No OpenTelemetry overhead** - Direct output, no batching or network calls

### 2. Automatic Trace Correlation

When logging within an active span (non-local environments), logs automatically include:

- `trace_id` - Links log to distributed trace
- `span_id` - Links log to specific span
- Service resource attributes

### 3. Configurable Log Levels

Different log levels for different environments:

```go
// Local development - see everything, pretty format
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

### 4. Unified Observability

All three signals (logs, traces, metrics) use the same:

- OTLP endpoint
- Resource attributes
- Correlation context

### 5. Production Ready

- **Batched processing** - Logs sent in efficient batches
- **Async delivery** - Non-blocking log processing
- **Resource attributes** - Proper service identification
- **Level filtering** - Efficient log processing
- **Error handling** - Graceful degradation

## Examples

See [examples/main.go](examples/main.go) for complete usage examples.
