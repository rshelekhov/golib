package metrics

import (
	"context"
	"net/http"
	"time"

	promclient "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type ExporterType string

const (
	ExporterPrometheus ExporterType = "prometheus"
	ExporterOTLP       ExporterType = "otlp"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	Env            string
	ExporterType   ExporterType
	OTLPEndpoint   string        // Used only when ExporterType is ExporterOTLP
	PushInterval   time.Duration // Used for OTLP exporter, defaults to 30s
	OTLPInsecure   bool          // If true, uses insecure OTLP connection
}

// Init initializes OpenTelemetry MeterProvider with the specified exporter
func Init(ctx context.Context, cfg Config) (*sdkmetric.MeterProvider, http.Handler, error) {
	// Create resource
	res := resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.ServiceVersion),
		semconv.DeploymentEnvironment(cfg.Env),
	)

	var provider *sdkmetric.MeterProvider
	var handler http.Handler
	var err error

	switch cfg.ExporterType {
	case ExporterOTLP:
		provider, err = initOTLP(ctx, res, cfg.OTLPEndpoint, cfg.PushInterval, cfg.OTLPInsecure)
	default: // ExporterPrometheus or empty
		provider, handler, err = initPrometheus(res)
	}

	if err != nil {
		return nil, nil, err
	}

	// Set global MeterProvider
	otel.SetMeterProvider(provider)

	return provider, handler, nil
}

func initPrometheus(res *resource.Resource) (*sdkmetric.MeterProvider, http.Handler, error) {
	// Create Prometheus exporter
	registry := promclient.NewRegistry()
	exporter, err := prometheus.New(
		prometheus.WithRegisterer(registry),
	)
	if err != nil {
		return nil, nil, err
	}

	// Create MeterProvider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)

	// Create HTTP handler for metrics
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	return provider, handler, nil
}

func initOTLP(ctx context.Context, res *resource.Resource, endpoint string, interval time.Duration, insecure bool) (*sdkmetric.MeterProvider, error) {
	// Create OTLP exporter with configurable TLS
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(endpoint),
	}
	if insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if interval == 0 {
		interval = 30 * time.Second
	}

	// Create periodic reader
	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(interval),
	)

	// Create MeterProvider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	return provider, nil
}

// OtelMeter returns global otel.Meter
func OtelMeter() metric.Meter {
	return otel.GetMeterProvider().Meter("github.com/rshelekhov/golib/observability/metrics")
}
