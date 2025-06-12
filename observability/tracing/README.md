# observability/tracing

OpenTelemetry tracing for Go microservices: simple initialization following observability course patterns.

## Quick start

### 1. Tracing initialization

```go
import (
    "context"
    "github.com/rshelekhov/golib/observability/tracing"
)

func main() {
    // Standard pattern from observability courses
    tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
    if err != nil {
        panic(err)
    }
    defer func() {
        if err := tracerProvider.Shutdown(context.Background()); err != nil {
            log.Printf("Error shutting down tracer: %v", err)
        }
    }()
    // ...
}
```

## Initialization Functions

### Stdout (Development)

```go
// Standard pattern - stdout exporter with propagation setup
tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
defer tracerProvider.Shutdown(ctx)

// Traces printed to stdout with pretty formatting
```

### OTLP (Production)

```go
// OTLP exporter for production
tracerProvider, err := tracing.InitTracerOTLP(ctx, "my-service", "1.0.0", "production", "localhost:4317")
defer tracerProvider.Shutdown(ctx)

// Traces sent to OTLP collector (Jaeger, Tempo, etc.)
```

## What's Configured Automatically

Both initialization functions set up:

1. **TracerProvider** with proper resource attributes
2. **Global TracerProvider** via `otel.SetTracerProvider()`
3. **TextMapPropagator** for trace context propagation:
   ```go
   otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
       propagation.TraceContext{},
       propagation.Baggage{},
   ))
   ```

## Resource Attributes

All functions automatically set:

```go
semconv.ServiceName(serviceName)
semconv.ServiceVersion(serviceVersion)
semconv.DeploymentEnvironment(env)
```

## Middleware for HTTP/gRPC

```go
import (
    "github.com/rshelekhov/golib/observability/tracing"
    "net/http"
    "google.golang.org/grpc"
)

// For HTTP
mux := http.NewServeMux()
handler := tracing.HTTPMiddleware(mux, "my-service")
log.Fatal(http.ListenAndServe(":8080", handler))

// For gRPC server
server := grpc.NewServer(
    grpc.StatsHandler(tracing.GRPCServerStatsHandler()),
)

// For gRPC client
conn, err := grpc.NewClient("localhost:50051",
    grpc.WithStatsHandler(tracing.GRPCClientStatsHandler()),
)
```

## Creating spans

```go
import "github.com/rshelekhov/golib/observability/tracing"

// In HTTP handler
ctx, span := tracing.SpanFromHTTP(r.Context(), r.Method, r.URL.Path)
defer span.End()

// In gRPC handler
ctx, span := tracing.SpanFromGRPC(ctx, "MyService.MyMethod")
defer span.End()

// For outgoing calls (DB, external services)
ctx, span := tracing.OutgoingSpan(ctx, "db.query", tracing.SpanKindClient,
    tracing.String("db.system", "postgresql"),
    tracing.String("db.statement", "SELECT * FROM users WHERE id = ?"),
)
defer span.End()

// Arbitrary span
ctx, span := tracing.StartSpan(ctx, "business.logic.operation")
defer span.End()
```

## Helpers for attributes

```go
// Instead of importing go.opentelemetry.io/otel/attribute
tracing.String("key", "value")
tracing.Int("count", 42)
tracing.Bool("success", true)

// Using in spans
span.SetAttributes(
    tracing.String("user.id", userID),
    tracing.Int("items.count", len(items)),
    tracing.Bool("cache.hit", true),
)
```

## Export to OTLP (Jaeger, Tempo, ...)

### Examples of endpoints:

- **Jaeger:** `localhost:14268` (HTTP) or `localhost:14250` (gRPC)
- **Tempo:** `localhost:4317` (gRPC) or `localhost:4318` (HTTP)
- **OTEL Collector:** `localhost:4317` (gRPC)

```go
// For production
tracerProvider, err := tracing.InitTracerOTLP(ctx, "my-service", "1.0.0", "production", "localhost:4317")

// For development
tracerProvider, err := tracing.InitTracer("my-service", "1.0.0", "development")
```

## Best Practices

- **Use `InitTracer()`** for development (stdout output)
- **Use `InitTracerOTLP()`** for production (OTLP collector)
- **Always call `defer tracerProvider.Shutdown(ctx)`**
- **Propagation is configured automatically** - no manual setup needed
- **Always close span** through `defer span.End()`
- **Pass context** between layers of the application
- **For HTTP/gRPC** always use middleware from tracing
- **For DB and external services** use `OutgoingSpan`

## Examples

- [examples/main.go](examples/main.go) - ready-made patterns for use
- [../examples/main.go](../examples/main.go) - full example with HTTP + gRPC
