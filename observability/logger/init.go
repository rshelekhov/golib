package logger

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitLogger initializes OpenTelemetry LoggerProvider with stdout exporter
func InitLogger(serviceName, serviceVersion, env string) (*log.LoggerProvider, *slog.Logger, error) {
	// Create stdout exporter
	exporter, err := stdoutlog.New()
	if err != nil {
		return nil, nil, fmt.Errorf("create stdout log exporter: %w", err)
	}

	// Create resource
	res := resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.DeploymentEnvironment(env),
	)

	// Create LoggerProvider
	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(res),
	)

	// Set global LoggerProvider
	global.SetLoggerProvider(lp)

	// Create OTel slog logger
	otelLogger := otelslog.NewLogger(serviceName)

	return lp, otelLogger, nil
}

// InitLoggerOTLP initializes OpenTelemetry LoggerProvider with OTLP exporter
func InitLoggerOTLP(ctx context.Context, serviceName, serviceVersion, env, endpoint string) (*log.LoggerProvider, *slog.Logger, error) {
	// Create OTLP exporter
	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create otlp log exporter: %w", err)
	}

	// Create resource
	res := resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.DeploymentEnvironment(env),
	)

	// Create LoggerProvider
	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(res),
	)

	// Set global LoggerProvider
	global.SetLoggerProvider(lp)

	// Create OTel slog logger
	otelLogger := otelslog.NewLogger(serviceName)

	return lp, otelLogger, nil
} 