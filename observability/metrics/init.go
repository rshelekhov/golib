package metrics

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	promclient "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitMeter initializes OpenTelemetry MeterProvider with Prometheus exporter
func InitMeter(serviceName, serviceVersion, env string) (*sdkmetric.MeterProvider, http.Handler, error) {
	// Create resource
	res := resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.DeploymentEnvironment(env),
	)

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

	// Set global MeterProvider
	otel.SetMeterProvider(provider)

	// Create HTTP handler for metrics
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	return provider, handler, nil
}

// InitMeterOTLP initializes OpenTelemetry MeterProvider with OTLP exporter
func InitMeterOTLP(ctx context.Context, serviceName, serviceVersion, env, endpoint string) (*sdkmetric.MeterProvider, error) {
	// Create resource
	res := resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.DeploymentEnvironment(env),
	)

	// Create OTLP exporter
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create periodic reader
	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(30*time.Second),
	)

	// Create MeterProvider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	// Set global MeterProvider
	otel.SetMeterProvider(provider)

	return provider, nil
}

// InitMeterStdout initializes OpenTelemetry MeterProvider with stdout exporter (for development)
func InitMeterStdout(serviceName, serviceVersion, env string) (*sdkmetric.MeterProvider, error) {
	// Create resource
	res := resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.DeploymentEnvironment(env),
	)

	// Create stdout exporter
	exporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	// Create periodic reader
	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(10*time.Second),
	)

	// Create MeterProvider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	// Set global MeterProvider
	otel.SetMeterProvider(provider)

	return provider, nil
}

// OtelMeter returns global otel.Meter
func OtelMeter() metric.Meter {
	return otel.GetMeterProvider().Meter("github.com/rshelekhov/golib/observability/metrics")
}
