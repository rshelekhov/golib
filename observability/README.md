# Observability

Library for logging, tracing and metrics in Go microservices and projects.

## Features

- **Logging** (slog, pretty colorized output for local, trace_id/span_id automatically, configurable log levels)
- **Tracing** (OpenTelemetry, stdout/OTLP exporters, propagation setup)
- **Metrics** (OpenTelemetry + Prometheus/OTLP exporters)
- **Configurable TLS** (secure/insecure OTLP connections with smart environment-based defaults)

## Quick start

### Local Development

```go
import (
	"context"
	"log"
	"net/http"

	"github.com/rshelekhov/golib/observability"
)

func main() {
	// Simple initialization for local development (metrics always disabled)
	cfg, err := observability.NewConfig(observability.ConfigParams{
		Env:            observability.EnvLocal,
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
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
		// Pretty colorized logs will be shown in local development
		obs.Logger.DebugContext(r.Context(), "processing request", "path", r.URL.Path)
		obs.Logger.InfoContext(r.Context(), "hello world")
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}))

	log.Printf("HTTP server listening on :8080")
	log.Printf("Metrics are disabled for local development")
	log.Printf("Logs use pretty colorized format for better readability")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Production

```go
func main() {
	// Production setup with default info logging and OTLP
	cfg, err := observability.NewConfig(observability.ConfigParams{
		Env:            observability.EnvProd,
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		EnableMetrics:  true,
		OTLPEndpoint:   "localhost:4317",
	})
	if err != nil {
		log.Fatal(err)
	}

	obs, err := observability.Init(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer obs.Shutdown(context.Background())

	// Debug logs won't be shown in production (info level)
	obs.Logger.DebugContext(ctx, "this won't be logged")
	obs.Logger.InfoContext(ctx, "this will be logged")
}
```

### Custom Log Levels

```go
// Override environment default log level
cfg, err := observability.NewConfig(
	observability.ConfigParams{
		Env:            observability.EnvLocal,
		ServiceName:    "k8s-service",
		ServiceVersion: "1.0.0",
		EnableMetrics:  false,
		OTLPEndpoint:   "",
	},
	observability.WithLogLevel(slog.LevelWarn), // Override default debug level to warn
)

obs, err := observability.Init(context.Background(), cfg)
```

### Error Handling for Invalid Configurations

```go
// This will return an error - unknown environment
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:            "staging",
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	EnableMetrics:  true,
	OTLPEndpoint:   "localhost:4317",
})
if err != nil {
	log.Fatal(err) // "unsupported environment: staging (supported: local, dev, prod)"
}

// This will return an error - missing OTLP endpoint for prod
cfg, err = observability.NewConfig(observability.ConfigParams{
	Env:            observability.EnvProd,
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	EnableMetrics:  true,
	OTLPEndpoint:   "", // Missing endpoint
})
if err != nil {
	log.Fatal(err) // "OTLP endpoint is required for environment prod"
}
```

### With Docker Observability Stack

For full observability experience with Jaeger, Prometheus, and Grafana:

```bash
# Start observability stack
cd infrastructure
docker-compose -f docker-compose.dev.yml up -d

# Access UIs:
# - Grafana: http://localhost:3000 (admin/admin)
# - Jaeger: http://localhost:16686
# - Prometheus: http://localhost:9090
```

Then use OTLP in your application:

```go
// Use OTLP to send all data to the stack
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:            observability.EnvDev,
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	EnableMetrics:  true,
	OTLPEndpoint:   "localhost:4317",
})
if err != nil {
	log.Fatal(err)
}
obs, err := observability.Init(context.Background(), cfg)
```

See [infrastructure/README.md](infrastructure/README.md) for setup guide, troubleshooting, and production considerations.

## TLS Configuration

The library provides configurable TLS support for OTLP connections with smart environment-based defaults:

### Default TLS Behavior

- **Local/Dev environments**: Insecure connections (no TLS) by default
- **Production environment**: Secure connections (TLS) by default

### Basic Usage

```go
// Production with TLS (default)
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:               observability.EnvProd,
	ServiceName:       "my-service",
	ServiceVersion:    "1.0.0",
	EnableMetrics:     true,
	OTLPEndpoint:      "otel-collector.company.com:4317",
	OTLPTransportType: tracing.OTLPGRPC,
	// OTLPInsecure: false (default for production - uses TLS)
})

// Development with insecure connection (default)
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:               observability.EnvDev,
	ServiceName:       "my-service",
	ServiceVersion:    "1.0.0",
	EnableMetrics:     true,
	OTLPEndpoint:      "localhost:4317",
	OTLPTransportType: tracing.OTLPGRPC,
	// OTLPInsecure: true (default for dev - no TLS)
})
```

### Override TLS Settings

Use functional options to override environment defaults:

```go
// Force TLS in development environment
cfg, err := observability.NewConfig(
	observability.ConfigParams{
		Env:               observability.EnvDev,
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "secure-collector.dev.company.com:4317",
		OTLPTransportType: tracing.OTLPGRPC,
	},
	observability.WithOTLPInsecure(false), // Override to use TLS
)

// Force insecure in production (not recommended)
cfg, err := observability.NewConfig(
	observability.ConfigParams{
		Env:               observability.EnvProd,
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "localhost:4317",
		OTLPTransportType: tracing.OTLPGRPC,
	},
	observability.WithOTLPInsecure(true), // Override to disable TLS
)
```

### Explicit Configuration

You can also set TLS configuration directly in ConfigParams:

```go
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:               observability.EnvProd,
	ServiceName:       "my-service",
	ServiceVersion:    "1.0.0",
	EnableMetrics:     true,
	OTLPEndpoint:      "localhost:4317",
	OTLPTransportType: tracing.OTLPGRPC,
	OTLPInsecure:      &[]bool{true}[0], // Explicitly disable TLS
})
```

### TLS Configuration Examples

```go
// Example 1: Local development with Docker Compose
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:               observability.EnvLocal,
	ServiceName:       "my-app",
	ServiceVersion:    "1.0.0",
	EnableMetrics:     false,
	OTLPEndpoint:      "", // No OTLP for local
	// TLS not relevant for local (no OTLP)
})

// Example 2: Development with local OTLP collector (insecure)
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:               observability.EnvDev,
	ServiceName:       "my-app",
	ServiceVersion:    "1.0.0",
	EnableMetrics:     true,
	OTLPEndpoint:      "localhost:4317",
	OTLPTransportType: tracing.OTLPGRPC,
	// OTLPInsecure: true (default for dev)
})

// Example 3: Staging with secure OTLP collector
cfg, err := observability.NewConfig(
	observability.ConfigParams{
		Env:               observability.EnvDev, // Using dev env
		ServiceName:       "my-app",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "otel-staging.company.com:4317",
		OTLPTransportType: tracing.OTLPGRPC,
	},
	observability.WithOTLPInsecure(false), // Override to use TLS
)

// Example 4: Production with TLS (default)
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:               observability.EnvProd,
	ServiceName:       "my-app",
	ServiceVersion:    "1.0.0",
	EnableMetrics:     true,
	OTLPEndpoint:      "otel-prod.company.com:4317",
	OTLPTransportType: tracing.OTLPGRPC,
	// OTLPInsecure: false (default for prod - uses TLS)
})
```

### Transport Types and TLS

TLS configuration works with both transport types:

- **GRPC Transport** (`tracing.OTLPGRPC`): Uses gRPC with configurable TLS
- **HTTP Transport** (`tracing.OTLPHTTP`): Uses HTTP with configurable TLS

```go
// GRPC with TLS
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:               observability.EnvProd,
	ServiceName:       "my-service",
	ServiceVersion:    "1.0.0",
	EnableMetrics:     true,
	OTLPEndpoint:      "otel-collector.company.com:4317",
	OTLPTransportType: tracing.OTLPGRPC, // Uses TLS by default in prod
})

// HTTP with TLS
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:               observability.EnvProd,
	ServiceName:       "my-service",
	ServiceVersion:    "1.0.0",
	EnableMetrics:     true,
	OTLPEndpoint:      "https://otel-collector.company.com:4318",
	OTLPTransportType: tracing.OTLPHTTP, // Uses TLS by default in prod
})
```

## Configuration Options

### Simplified Configuration API

```go
// Using ConfigParams struct with optional functional options
// ConfigParams provides type safety and clear parameter names

// Local development with default debug logging (metrics always disabled)
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:            observability.EnvLocal,
	ServiceName:    "service-name",
	ServiceVersion: "1.0.0",
	EnableMetrics:  false,
	OTLPEndpoint:   "",
})

// Production with default info logging
cfg, err := observability.NewConfig(observability.ConfigParams{
	Env:            observability.EnvProd,
	ServiceName:    "service-name",
	ServiceVersion: "1.0.0",
	EnableMetrics:  true,
	OTLPEndpoint:   "localhost:4317",
})

// Override log level for any environment using functional options
cfg, err := observability.NewConfig(
	observability.ConfigParams{
		Env:            observability.EnvLocal,
		ServiceName:    "service-name",
		ServiceVersion: "1.0.0",
		EnableMetrics:  false,
		OTLPEndpoint:   "",
	},
	observability.WithLogLevel(slog.LevelWarn),
)

cfg, err := observability.NewConfig(
	observability.ConfigParams{
		Env:            observability.EnvProd,
		ServiceName:    "service-name",
		ServiceVersion: "1.0.0",
		EnableMetrics:  true,
		OTLPEndpoint:   "endpoint",
	},
	observability.WithLogLevel(slog.LevelError),
)
```

### Environment Defaults

- **Local**: Debug log level, pretty colorized output, no OTLP endpoint required
- **Dev/Prod**: Info log level, OTLP output, endpoint required

### Manual Configuration

```go
// Direct struct initialization still works for advanced cases
cfg := observability.Config{
	Env:            "local",
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	EnableMetrics:  true,
	OTLPEndpoint:   "",
	LogLevel:       slog.LevelDebug,
}

// Manual configs don't get automatic validation - need to validate parameters manually
params := observability.ConfigParams{
	Env:            cfg.Env,
	ServiceName:    cfg.ServiceName,
	ServiceVersion: cfg.ServiceVersion,
	EnableMetrics:  cfg.EnableMetrics,
	OTLPEndpoint:   cfg.OTLPEndpoint,
}
if err := params.Validate(); err != nil {
	log.Fatal(err)
}
```

### Log Levels

- `slog.LevelDebug` - All logs (local default)
- `slog.LevelInfo` - Info and above (dev/prod default)
- `slog.LevelWarn` - Warnings and errors only
- `slog.LevelError` - Errors only

## Automatic Exporter Selection

The `Init()` function automatically chooses the appropriate exporters based on your configuration:

### Local Development (`EnvLocal`)

- **Logging**: pretty colorized handler with human-readable format
- **Tracing**: stdout with pretty formatting
- **Metrics**: completely disabled (no overhead)

### Production (`EnvProd`, `EnvDev` with OTLP endpoint)

- **Logging**: OTLP exporter
- **Tracing**: OTLP exporter
- **Metrics**: OTLP exporter (push model)

### Configuration Examples

```go
// Local development - uses pretty logging, metrics always disabled
cfg := observability.Config{
	Env:            observability.EnvLocal,
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	EnableMetrics:  false, // Ignored for local - always disabled
	// OTLPEndpoint not needed for local
}

// Production - uses OTLP for everything
cfg := observability.Config{
	Env:            observability.EnvProd,
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	EnableMetrics:  true,
	OTLPEndpoint:   "localhost:4317", // Required for prod
}

obs, err := observability.Init(context.Background(), cfg)
```

## Architecture

All sub-packages follow the same pattern:

- **`logger.Init(ctx, logger.Config)`** - unified logger initialization
- **`metrics.Init(ctx, metrics.Config)`** - unified metrics initialization
- **`tracing.Init(ctx, tracing.Config)`** - unified tracing initialization
- **`observability.Init(ctx, observability.Config)`** - orchestrates all components

Each package automatically selects the appropriate exporter based on configuration.

## Examples

See [examples/main.go](examples/main.go) for complete working examples of:

- Local development setup
- Production setup
- Custom log levels
- Error handling
- Manual component initialization

## Best practices

- **Use ConfigParams struct** - type-safe configuration with clear parameter names
- **Use functional options** - override environment defaults when needed
- **Use `observability.Init()`** - automatically chooses stdout or OTLP based on env
- **Adjust log levels** - Debug for local, Info for production, Warn for high-traffic services
- **Always call `defer obs.Shutdown(ctx)`** for proper cleanup
- **Propagation is set up automatically** - no manual configuration needed

## Integration with DB

- [db/postgres/README.md](../db/postgres/README.md) - PostgreSQL with tracing
- [db/mongo/README.md](../db/mongo/README.md) - MongoDB with tracing
