# Observability

Library for logging, tracing and metrics in Go microservices and projects.

## Features

- **Logging** (slog, trace_id/span_id automatically)
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
	// Simple initialization with stdout tracing and Prometheus metrics
	obs, err := observability.Init(context.Background(), observability.Config{
		Env:            "development",
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		EnableMetrics:  true, // Uses Prometheus by default
	})
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
		obs.Logger.InfoContext(r.Context(), "hello world")
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})))

	log.Printf("HTTP server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
obs, err := observability.InitWithOTLP(context.Background(), observability.Config{
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	Env:            "development",
	EnableMetrics:  true,
}, "localhost:4317") // OTel Collector endpoint
```

See [infrastructure/README.md](infrastructure/README.md) for setup guide, troubleshooting, and production considerations.

## Initialization Options

### Simple (Development)

```go
// Stdout tracing + Prometheus metrics
obs, err := observability.Init(ctx, observability.Config{
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	Env:            "development",
	EnableMetrics:  true,
})
```

### Production (OTLP)

```go
// OTLP tracing + OTLP metrics
obs, err := observability.InitWithOTLP(ctx, observability.Config{
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	Env:            "production",
	EnableMetrics:  true,
}, "localhost:4317") // OTLP endpoint
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
- [examples/main.go](examples/main.go) - full example

## Best practices

- **Use simple `observability.Init()`** for development (stdout tracing + Prometheus)
- **Use `InitWithOTLP()`** for production (OTLP tracing + metrics)
- **Use direct initialization** when you need fine control
- **Always call `defer obs.Shutdown(ctx)`** for proper cleanup
- **Propagation is set up automatically** - no manual configuration needed

## Integration with DB

- [db/postgres/README.md](../db/postgres/README.md) - PostgreSQL with tracing
- [db/mongo/README.md](../db/mongo/README.md) - MongoDB with tracing
