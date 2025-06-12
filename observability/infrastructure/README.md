# Observability Stack with Docker

Docker-based observability stack for development and testing.

## Files Overview

### Docker Compose

- `docker-compose.dev.yml` - Simplified stack for development (recommended)
- `docker-compose.yml` - Full production-like stack with Loki, Tempo

### OpenTelemetry Collector

- `otel-collector-dev.yml` - Simple config for development
- `otel-collector-config.yml` - Full config with all exporters

### Prometheus

- `prometheus-dev.yml` - Basic Prometheus config
- `prometheus.yml` - Full config with multiple targets

### Grafana

- `grafana-dev.yml` - Auto-configured datasources

## Quick Start (Development)

### 1. Start the Stack

```bash
# Start simplified stack for development
docker-compose -f docker-compose.dev.yml up -d

# Check status
docker-compose -f docker-compose.dev.yml ps
```

### 2. Configure Your Application

In your Go application:

```go
// Use OTLP endpoint
obs, err := observability.InitWithOTLP(context.Background(), observability.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:           "development",
    EnableMetrics: true,
}, "localhost:4317")
```

### 3. Access UIs

- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686
- **Prometheus**: http://localhost:9090

## Full Stack (Production-like)

```bash
# Start full stack
docker-compose up -d

# Stop
docker-compose down

# Stop with data removal
docker-compose down -v
```

## Stack Components

### OpenTelemetry Collector

- **Port**: 4317 (gRPC), 4318 (HTTP)
- **Purpose**: Central collection point for all observability data
- **Configuration**: `otel-collector-config.yml`

### Jaeger

- **UI**: http://localhost:16686
- **Purpose**: Distributed tracing, trace visualization
- **Data**: Traces from your application

### Prometheus

- **UI**: http://localhost:9090
- **Purpose**: Metrics collection and storage
- **Configuration**: `prometheus.yml`

### Grafana

- **UI**: http://localhost:3000 (admin/admin)
- **Purpose**: Unified interface for all data
- **Datasources**: Automatically configured Prometheus and Jaeger

### Loki (full stack only)

- **Port**: 3100
- **Purpose**: Log collection and storage
- **Integration**: Through OTel Collector

## Integration with Your Application

### 1. OTLP (recommended)

```go
package main

import (
    "context"
    "net/http"

    "github.com/rshelekhov/golib/observability"
    "github.com/rshelekhov/golib/observability/metrics"
)

func main() {
    // Initialize with OTLP
    obs, err := observability.InitWithOTLP(context.Background(),
        observability.Config{
            ServiceName:    "my-service",
            ServiceVersion: "1.0.0",
            Env:           "development",
            EnableMetrics: true,
        },
        "localhost:4317", // OTel Collector endpoint
    )
    if err != nil {
        panic(err)
    }
    defer obs.Shutdown(context.Background())

    // HTTP server with metrics and tracing
    http.Handle("/", metrics.Middleware(http.HandlerFunc(handler)))
    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Logs automatically correlate with traces
    obs.Logger.InfoContext(r.Context(), "handling request",
        "path", r.URL.Path,
        "method", r.Method,
    )

    w.Write([]byte("Hello World!"))
}
```

### 2. Prometheus + stdout (alternative)

```go
// For cases when you need Prometheus endpoint
obs, err := observability.Init(context.Background(), observability.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Env:           "development",
    EnableMetrics: true,
})

// Add Prometheus endpoint
http.Handle("/metrics", obs.MetricsHandler)
```

## Monitoring and Debugging

### Check Collector Status

```bash
# Collector logs
docker logs otel-collector-dev

# Should see messages about receiving data:
# 2024-01-20T10:30:45.123Z info TracesExporter {"#traces": 1}
# 2024-01-20T10:30:45.124Z info MetricsExporter {"#metrics": 5}
```

### Check Metrics in Prometheus

1. Open http://localhost:9090
2. In query field enter: `{__name__=~".*my_service.*"}`
3. Your service metrics should appear

### Check Traces in Jaeger

1. Open http://localhost:16686
2. Select service "my-service"
3. Click "Find Traces"

## Troubleshooting

### Application Not Sending Data

```bash
# Check if collector is accessible
curl http://localhost:4317

# Check application logs for OTLP errors
```

### Metrics Not Appearing in Prometheus

```bash
# Check targets in Prometheus
# http://localhost:9090/targets

# Check collector configuration
docker logs otel-collector-dev | grep -i error
```

### Traces Not Visible in Jaeger

```bash
# Check if Jaeger is receiving data
docker logs jaeger-dev | grep -i trace

# Check pipeline in collector config
```

## Customization

### Adding Alerts

Create `alert_rules.yml`:

```yaml
groups:
  - name: my-service-alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 2m
        annotations:
          summary: "High error rate detected"
```

### Adding Dashboards

Place JSON dashboards in `grafana/dashboards/` and they will be automatically loaded.

## Production Considerations

1. **Persistent storage**: Add volumes for data
2. **Security**: Configure authentication and TLS
3. **Scaling**: Use multiple collectors
4. **Retention**: Configure data retention policies
5. **Alerting**: Add Alertmanager

## Useful Commands

```bash
# Restart only collector
docker-compose restart otel-collector

# View logs of all services
docker-compose logs -f

# Clean all data
docker-compose down -v
docker system prune -f

# Update images
docker-compose pull
docker-compose up -d
```
