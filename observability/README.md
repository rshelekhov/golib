# Observability

Library for logging, tracing and metrics in Go microservices and projects.

## Features

- **Logging** (slog, trace_id/span_id automatically, configurable log levels)
- **Tracing** (OpenTelemetry, stdout/OTLP exporters, propagation setup)
- **Metrics** (OpenTelemetry + Prometheus/OTLP exporters)

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
	cfg, err := observability.NewConfig(observability.EnvLocal, "my-service", "1.0.0", false, "")
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
		// Debug logs will be shown in local development
		obs.Logger.DebugContext(r.Context(), "processing request", "path", r.URL.Path)
		obs.Logger.InfoContext(r.Context(), "hello world")
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}))

	log.Printf("HTTP server listening on :8080")
	log.Printf("Metrics are disabled for local development")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Production

```go
func main() {
	// Production setup with default info logging and OTLP
	cfg, err := observability.NewConfig(observability.EnvProd, "my-service", "1.0.0", true, "localhost:4317")
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
	observability.EnvLocal, "k8s-service", "1.0.0", true, "",
	slog.LevelWarn, // Override default debug level to warn
)

obs, err := observability.Init(context.Background(), cfg)
```

### Error Handling for Invalid Configurations

```go
// This will return an error - unknown environment
cfg, err := observability.NewConfig("staging", "my-service", "1.0.0", true, "")
if err != nil {
	log.Fatal(err) // "unsupported environment: staging (supported: local, dev, prod)"
}

// This will return an error - missing OTLP endpoint for prod
cfg, err = observability.NewConfig(observability.EnvProd, "my-service", "1.0.0", true, "")
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
cfg, err := observability.NewConfig(observability.EnvDev, "my-service", "1.0.0", true, "localhost:4317")
if err != nil {
	log.Fatal(err)
}
obs, err := observability.Init(context.Background(), cfg)
```

See [infrastructure/README.md](infrastructure/README.md) for setup guide, troubleshooting, and production considerations.

## Configuration Options

### Simplified Configuration API

```go
// NewConfig(env, serviceName, serviceVersion, enableMetrics, otlpEndpoint, logLevel...)
// Optional logLevel parameter overrides environment defaults

// Local development with default debug logging (metrics always disabled)
cfg, err := observability.NewConfig(observability.EnvLocal, "service-name", "1.0.0", false, "")

// Production with default info logging
cfg, err := observability.NewConfig(observability.EnvProd, "service-name", "1.0.0", true, "localhost:4317")

// Override log level for any environment
cfg, err := observability.NewConfig(observability.EnvLocal, "service-name", "1.0.0", true, "", slog.LevelWarn)
cfg, err := observability.NewConfig(observability.EnvProd, "service-name", "1.0.0", true, "endpoint", slog.LevelError)
```

### Environment Defaults

- **Local**: Debug log level, stdout output, no OTLP endpoint required
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

// Always validate manually created configs
if err := cfg.Validate(); err != nil {
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

- **Logging**: stdout with pretty formatting
- **Tracing**: stdout with pretty formatting
- **Metrics**: completely disabled (no overhead)

### Production (`EnvProd`, `EnvDev` with OTLP endpoint)

- **Logging**: OTLP exporter
- **Tracing**: OTLP exporter
- **Metrics**: OTLP exporter (push model)

### Configuration Examples

```go
// Local development - uses stdout, metrics always disabled
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

- **Use helper functions** - `NewLocalConfig()`, `NewProdConfig()` for typical setups
- **Use `observability.Init()`** - automatically chooses stdout or OTLP based on env
- **Adjust log levels** - Debug for local, Info for production, Warn for high-traffic services
- **Always call `defer obs.Shutdown(ctx)`** for proper cleanup
- **Propagation is set up automatically** - no manual configuration needed

## Integration with DB

- [db/postgres/README.md](../db/postgres/README.md) - PostgreSQL with tracing
- [db/mongo/README.md](../db/mongo/README.md) - MongoDB with tracing
