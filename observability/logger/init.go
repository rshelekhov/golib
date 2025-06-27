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
func InitLoggerStdout(serviceName, serviceVersion, env string, level slog.Level) (*log.LoggerProvider, *slog.Logger, error) {
	// Create stdout exporter
	exporter, err := stdoutlog.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stdout log exporter: %w", err)
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

	// Create OTEL handler with level filtering
	handler := otelslog.NewHandler(serviceName, otelslog.WithLoggerProvider(lp))
	finalLogger := slog.New(&levelFilterHandler{
		handler:  handler,
		minLevel: level,
	})

	return lp, finalLogger, nil
}

// InitLoggerOTLP initializes OpenTelemetry LoggerProvider with OTLP exporter
func InitLoggerOTLP(ctx context.Context, serviceName, serviceVersion, env, endpoint string, level slog.Level) (*log.LoggerProvider, *slog.Logger, error) {
	// Create OTLP exporter
	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create otlp log exporter: %w", err)
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

	// Create slog logger with level filtering
	handler := otelslog.NewHandler(serviceName, otelslog.WithLoggerProvider(lp))
	finalLogger := slog.New(&levelFilterHandler{
		handler:  handler,
		minLevel: level,
	})

	return lp, finalLogger, nil
}

// levelFilterHandler wraps a slog.Handler to filter by log level
type levelFilterHandler struct {
	handler  slog.Handler
	minLevel slog.Level
}

func (h *levelFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.minLevel && h.handler.Enabled(ctx, level)
}

func (h *levelFilterHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= h.minLevel {
		return h.handler.Handle(ctx, record)
	}
	return nil
}

func (h *levelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelFilterHandler{
		handler:  h.handler.WithAttrs(attrs),
		minLevel: h.minLevel,
	}
}

func (h *levelFilterHandler) WithGroup(name string) slog.Handler {
	return &levelFilterHandler{
		handler:  h.handler.WithGroup(name),
		minLevel: h.minLevel,
	}
}
