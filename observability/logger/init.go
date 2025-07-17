package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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

type Config struct {
	ServiceName    string
	ServiceVersion string
	Env            string
	Level          slog.Level
	Endpoint       string // OTLP endpoint. If empty, stdout exporter is used.
	OTLPInsecure   bool   // If true, uses insecure OTLP connection
}

// Init initializes OpenTelemetry LoggerProvider
func Init(ctx context.Context, cfg Config) (*log.LoggerProvider, *slog.Logger, error) {
	// For local environment, use pretty handler instead of OTEL
	if cfg.Env == "local" {
		handler := NewPrettyHandler(os.Stdout, &PrettyHandlerOptions{
			Level:     cfg.Level,
			AddSource: true,
		})

		finalLogger := slog.New(&levelFilterHandler{
			handler:  handler,
			minLevel: cfg.Level,
		})

		// Return nil LoggerProvider for local env since we're not using OTEL
		return nil, finalLogger, nil
	}

	var exporter log.Exporter
	var err error

	if cfg.Endpoint == "" {
		// Create stdout exporter
		exporter, err = stdoutlog.New()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create stdout log exporter: %w", err)
		}
	} else {
		// Create OTLP exporter with configurable TLS
		opts := []otlploggrpc.Option{
			otlploggrpc.WithEndpoint(cfg.Endpoint),
		}
		if cfg.OTLPInsecure {
			opts = append(opts, otlploggrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
		}

		exporter, err = otlploggrpc.New(ctx, opts...)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create otlp log exporter: %w", err)
		}
	}

	// Create resource
	res := resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.ServiceVersion),
		semconv.DeploymentEnvironment(cfg.Env),
	)

	// Create LoggerProvider
	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(res),
	)

	// Set global LoggerProvider
	global.SetLoggerProvider(lp)

	// Create slog logger with level filtering
	handler := otelslog.NewHandler(cfg.ServiceName, otelslog.WithLoggerProvider(lp))
	finalLogger := slog.New(&levelFilterHandler{
		handler:  handler,
		minLevel: cfg.Level,
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
