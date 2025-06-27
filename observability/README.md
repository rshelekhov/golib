# Observability

Library for logging, tracing and metrics in Go microservices and projects.

## Features

- **Logging** (slog, trace_id/span_id automatically, configurable log levels)
- **Tracing** (OpenTelemetry, stdout/OTLP exporters, propagation setup)
- **Metrics** (OpenTelemetry + Prometheus/OTLP/stdout exporters)

## Quick start

### Local Development

```go
import (
	"context"
	"log"
	"net/http"

	"github.com/rshelekhov/golib/observability"
	"github.com/rshelekhov/golib/observability/metrics"
)

func main() {
	// Simple initialization with default debug logging for local development
	cfg, err := observability.NewConfig(observability.EnvLocal, "my-service", "1.0.0", true, "")
	if err != nil {
		log.Fatal(err)
	}

	obs, err := observability.Setup(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := obs.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	// HTTP server with metrics endpoint
	http.Handle("/metrics", obs.MetricsHandler) // Prometheus scrapes this
	http.Handle("/", metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Debug logs will be shown in local development
		obs.Logger.DebugContext(r.Context(), "processing request", "path", r.URL.Path)
		obs.Logger.InfoContext(r.Context(), "hello world")
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})))

	log.Printf("HTTP server listening on :8080")
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

	obs, err := observability.Setup(context.Background(), cfg)
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

obs, err := observability.Setup(context.Background(), cfg)
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
obs, err := observability.Setup(context.Background(), cfg)
```

See [infrastructure/README.md](infrastructure/README.md) for setup guide, troubleshooting, and production considerations.

## Configuration Options

### Simplified Configuration API

```go
// NewConfig(env, serviceName, serviceVersion, enableMetrics, otlpEndpoint, logLevel...)
// Optional logLevel parameter overrides environment defaults

// Local development with default debug logging
cfg, err := observability.NewConfig(observability.EnvLocal, "service-name", "1.0.0", true, "")

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

## Initialization Options

### Simple (Automatic based on Env)

```go
cfg := observability.NewLocalConfig("my-service", "1.0.0", true)
obs, err := observability.Setup(context.Background(), cfg)
```

### Advanced (Custom)

```go
// Custom OTLP endpoint
obs, err := observability.InitWithOTLP(ctx, cfg, "custom-endpoint:4317")
```

## Direct Initialization (Advanced)

### Tracing

```go
// Initialize tracer directly (like in observability courses)
tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
defer tracerProvider.Shutdown(ctx)

// Or OTLP tracer
tracerProvider, err := tracing.InitTracerOTLP(ctx, "my-service", "1.0.0", "production", "localhost:4317")
defer tracerProvider.Shutdown(ctx)
```

### Logging

```go
// Direct logger initialization with custom level
loggerProvider, logger, err := logger.InitLoggerStdout("my-service", "1.0.0", "development", slog.LevelWarn)
defer loggerProvider.Shutdown(ctx)

// Or OTLP logger
loggerProvider, logger, err := logger.InitLoggerOTLP(ctx, "my-service", "1.0.0", "production", "localhost:4317", slog.LevelInfo)
defer loggerProvider.Shutdown(ctx)
```

### Metrics

```go
// Initialize meter directly
meterProvider, handler, err := metrics.InitMeter("my-service", "1.0.0", "production")
defer meterProvider.Shutdown(ctx)
http.Handle("/metrics", handler)
```

## Features

### Automatic Propagation Setup

Both `InitTracer` and `InitTracerOTLP` automatically configure:

```go
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
	propagation.TraceContext{},
	propagation.Baggage{},
))
```

### Resource Attributes

All initialization functions set:

```go
semconv.ServiceName(serviceName)
semconv.ServiceVersion(serviceVersion)
semconv.DeploymentEnvironment(env)
```

## gRPC integration

```go
// Use stats.Handler for tracing
server := grpc.NewServer(
	grpc.StatsHandler(tracing.GRPCServerStatsHandler()),
	grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()),
	grpc.StreamInterceptor(metrics.StreamServerInterceptor()),
)

// For client
conn, err := grpc.NewClient("localhost:50051",
	grpc.WithStatsHandler(tracing.GRPCClientStatsHandler()),
)
```

## Creating spans manually

```go
// HTTP span
ctx, span := tracing.SpanFromHTTP(ctx, "GET", "/api/v1/users/{id}")
defer span.End()

// gRPC span
ctx, span := tracing.SpanFromGRPC(ctx, "UserService.GetUser")
defer span.End()

// Outgoing call (DB, external API)
ctx, span := tracing.OutgoingSpan(ctx, "db.query", tracing.SpanKindClient,
	tracing.String("db.system", "postgresql"),
	tracing.String("db.statement", "SELECT * FROM users WHERE id = ?"),
)
defer span.End()
```

## Business metrics

```go
// Write business error
metrics.IncBusinessError("validation", "invalid_email")
```

## Examples

- [logger/README.md](logger/README.md)
- [tracing/README.md](tracing/README.md)
- [metrics/README.md](metrics/README.md)
- [examples/main.go](examples/main.go) - full example with different log levels

## Best practices

- **Use helper functions** - `NewLocalConfig()`, `NewProdConfig()` for typical setups
- **Use `observability.Setup()`** - automatically chooses stdout or OTLP based on env
- **Adjust log levels** - Debug for local, Info for production, Warn for high-traffic services
- **Always call `defer obs.Shutdown(ctx)`** for proper cleanup
- **Propagation is set up automatically** - no manual configuration needed

## Integration with DB

- [db/postgres/README.md](../db/postgres/README.md) - PostgreSQL with tracing
- [db/mongo/README.md](../db/mongo/README.md) - MongoDB with tracing
