# observability/tracing

OpenTelemetry tracing for Go microservices: simple initialization with automatic exporter selection.

## Quick start

### 1. Tracing initialization

```go
import (
    "context"
    "github.com/rshelekhov/golib/observability/tracing"
)

func main() {
    // Initialize with automatic exporter selection
    cfg := tracing.Config{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        Env:            "development",
        ExporterType:   tracing.ExporterStdout, // or tracing.ExporterOTLP
    }
    tracerProvider, err := tracing.Init(context.Background(), cfg)
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

## Configuration

```go
type Config struct {
    ServiceName    string
    ServiceVersion string
    Env            string
    ExporterType   ExporterType
    OTLPEndpoint   string // Used only when ExporterType is ExporterOTLP
}

type ExporterType string

const (
    ExporterStdout ExporterType = "stdout"
    ExporterOTLP   ExporterType = "otlp"
)
```

## Initialization

### `Init(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error)`

Unified initialization function that automatically selects the appropriate exporter.

**Exporter Selection:**

- `ExporterStdout`: Pretty-printed traces to stdout (development)
- `ExporterOTLP`: Traces sent to OTLP collector (production)

**Returns:**

- `*sdktrace.TracerProvider` - For shutdown management
- `error` - Initialization error

## Exporter Types

### Stdout (Development)

```go
// Standard pattern - stdout exporter with propagation setup
cfg := tracing.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:            "development",
    ExporterType:   tracing.ExporterStdout,
}
tracerProvider, err := tracing.Init(ctx, cfg)
defer tracerProvider.Shutdown(ctx)

// Traces printed to stdout with pretty formatting
```

### OTLP (Production)

```go
// OTLP exporter for production
cfg := tracing.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:            "production",
    ExporterType:   tracing.ExporterOTLP,
    OTLPEndpoint:   "localhost:4317",
}
tracerProvider, err := tracing.Init(ctx, cfg)
defer tracerProvider.Shutdown(ctx)

// Traces sent to OTLP collector (Jaeger, Tempo, etc.)
```

## What's Configured Automatically

The initialization function sets up:

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

All initialization automatically sets:

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
cfg := tracing.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:            "production",
    ExporterType:   tracing.ExporterOTLP,
    OTLPEndpoint:   "localhost:4317",
}
tracerProvider, err := tracing.Init(ctx, cfg)

// For development
cfg = tracing.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:            "development",
    ExporterType:   tracing.ExporterStdout,
}
tracerProvider, err = tracing.Init(ctx, cfg)
```

## Best Practices

- **Use `ExporterStdout`** for development (stdout output)
- **Use `ExporterOTLP`** for production (OTLP collector)
- **Always call `defer tracerProvider.Shutdown(ctx)`**
- **Propagation is configured automatically** - no manual setup needed
- **Always close span** through `defer span.End()`
- **Pass context** between layers of the application
- **For HTTP/gRPC** always use middleware from tracing
- **For DB and external services** use `OutgoingSpan`

## Examples

- [examples/main.go](examples/main.go) - ready-made patterns for use
- [../examples/main.go](../examples/main.go) - full example with HTTP + gRPC
