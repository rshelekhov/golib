# observability/metrics

OpenTelemetry metrics for Go services: simple initialization with automatic exporter selection.

## Quick start

```go
import (
	"context"
	"net/http"
	"github.com/rshelekhov/golib/observability/metrics"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	// Initialize meter with automatic exporter selection
	cfg := metrics.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "production",
		ExporterType:   metrics.ExporterPrometheus, // or metrics.ExporterOTLP
	}
	meterProvider, handler, err := metrics.Init(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down metrics: %v", err)
		}
	}()

	// Metrics endpoint (only for Prometheus exporter)
	if handler != nil {
		http.Handle("/metrics", handler)
	}

	// Create custom metric
	meter := metrics.OtelMeter()
	counter, err := meter.Int64Counter(
		"my_otel_counter",
		metric.WithDescription("Example otel counter."),
	)
	if err != nil {
		panic(err)
	}

	// Use metric
	counter.Add(context.Background(), 1)

	// HTTP middleware with metrics
	http.Handle("/", metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Configuration

```go
type Config struct {
	ServiceName    string
	ServiceVersion string
	Env            string
	ExporterType   ExporterType
	OTLPEndpoint   string        // Used only when ExporterType is ExporterOTLP
	PushInterval   time.Duration // Used for OTLP exporter, defaults to 30s
}

type ExporterType string

const (
	ExporterPrometheus ExporterType = "prometheus"
	ExporterOTLP       ExporterType = "otlp"
)
```

## Initialization

### `Init(ctx context.Context, cfg Config) (*sdkmetric.MeterProvider, http.Handler, error)`

Unified initialization function that automatically selects the appropriate exporter.

**Exporter Selection:**

- `ExporterPrometheus`: Pull model with HTTP endpoint for scraping
- `ExporterOTLP`: Push model, metrics sent to OTLP collector

**Returns:**

- `*sdkmetric.MeterProvider` - For shutdown management
- `http.Handler` - Metrics endpoint (only for Prometheus, nil for OTLP)
- `error` - Initialization error

## Exporter Types

### Prometheus (Pull Model)

```go
// Standard pattern for Prometheus scraping
cfg := metrics.Config{
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	Env:            "production",
	ExporterType:   metrics.ExporterPrometheus,
}
meterProvider, handler, err := metrics.Init(ctx, cfg)
defer meterProvider.Shutdown(ctx)

// Expose metrics via HTTP endpoint
http.Handle("/metrics", handler)
```

### OTLP (Push Model)

```go
// For Grafana Cloud, Tempo, or OTLP collector
cfg := metrics.Config{
	ServiceName:    "my-service",
	ServiceVersion: "1.0.0",
	Env:            "production",
	ExporterType:   metrics.ExporterOTLP,
	OTLPEndpoint:   "localhost:4317",
	PushInterval:   30 * time.Second, // Optional, defaults to 30s
}
meterProvider, _, err := metrics.Init(ctx, cfg)
defer meterProvider.Shutdown(ctx)

// Metrics automatically pushed every 30 seconds
```

## HTTP metrics

```go
// Automatic metrics for HTTP
handler := metrics.Middleware(yourHandler)
```

**Metrics:**

- `http_requests_total` - total number of HTTP requests
- `http_request_duration_seconds` - HTTP request processing time
- `http_panics_total` - number of panics in HTTP handlers

## gRPC metrics

```go
server := grpc.NewServer(
	grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()),
	grpc.StreamInterceptor(metrics.StreamServerInterceptor()),
)
```

**Metrics:**

- `grpc_server_requests_total` - total number of gRPC requests
- `grpc_server_handling_seconds` - gRPC request processing time

## Business errors

```go
// Record business error with type and code
metrics.IncBusinessError("validation", "empty_email")
metrics.IncBusinessError("auth", "invalid_token")
metrics.IncBusinessError("database", "connection_timeout")
```

**Metric:**

- `business_errors_total` - number of business errors with labels `type` and `code`

## Custom metrics

```go
meter := metrics.OtelMeter()

// Counter
counter, err := meter.Int64Counter(
	"orders_processed_total",
	metric.WithDescription("Total number of processed orders"),
)
counter.Add(ctx, 1, metric.WithAttributes(
	attribute.String("status", "completed"),
	attribute.String("region", "us-east"),
))

// Histogram
histogram, err := meter.Float64Histogram(
	"order_processing_duration_seconds",
	metric.WithDescription("Order processing duration in seconds"),
)
histogram.Record(ctx, duration.Seconds())

// Gauge (via UpDownCounter)
gauge, err := meter.Int64UpDownCounter(
	"active_connections",
	metric.WithDescription("Number of active connections"),
)
gauge.Add(ctx, 1)  // increase
gauge.Add(ctx, -1) // decrease
```

## Resource Attributes

All initialization functions automatically set:

```go
semconv.ServiceName(serviceName)
semconv.ServiceVersion(serviceVersion)
semconv.DeploymentEnvironment(env)
```

## Best practices

- **Use `ExporterPrometheus` for development and pull-based monitoring**
- **Use `ExporterOTLP` for production and push-based monitoring**
- **Always call `defer meterProvider.Shutdown(ctx)`**
- **Don't abuse label cardinality** (avoid user_id, ip as labels)
- **Use fixed values for labels** (error type, HTTP status)
- **Don't create metrics dynamically** in loops

## Available metrics

### HTTP metrics

- `http_requests_total{method, path, status}` - HTTP request counter
- `http_request_duration_seconds{method, path}` - processing time histogram
- `http_panics_total{method, path}` - panic counter

### gRPC metrics

- `grpc_server_requests_total{service, method, code}` - gRPC request counter
- `grpc_server_handling_seconds{service, method}` - processing time histogram

### Business metrics

- `business_errors_total{type, code}` - business error counter

## Examples

- [examples/main.go](examples/main.go) - ready-made usage patterns
- [../examples/main.go](../examples/main.go) - full example with HTTP + gRPC
