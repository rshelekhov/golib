package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExporterType string

const (
	ExporterStdout ExporterType = "stdout"
	ExporterOTLP   ExporterType = "otlp"
)

type OTLPTransportType string

const (
	OTLPTransportGRPC OTLPTransportType = "grpc"
	OTLPTransportHTTP OTLPTransportType = "http"
)

type Config struct {
	ServiceName       string
	ServiceVersion    string
	Env               string
	ExporterType      ExporterType
	OTLPEndpoint      string            // Used only when ExporterType is ExporterOTLP
	OTLPTransportType OTLPTransportType // "grpc" or "http", used only when ExporterType is ExporterOTLP
	OTLPInsecure      bool              // If true, uses insecure OTLP connection
}

// Init initializes OpenTelemetry TracerProvider
func Init(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error) {
	var exporter sdktrace.SpanExporter
	var err error

	switch cfg.ExporterType {
	case ExporterOTLP:
		switch cfg.OTLPTransportType {
		case OTLPTransportHTTP:
			opts := []otlptracehttp.Option{
				otlptracehttp.WithEndpoint(cfg.OTLPEndpoint),
			}
			if cfg.OTLPInsecure {
				opts = append(opts, otlptracehttp.WithInsecure())
			}

			exporter, err = otlptracehttp.New(ctx, opts...)
			if err != nil {
				return nil, fmt.Errorf("create otlp http exporter: %w", err)
			}
		case OTLPTransportGRPC:
			opts := []otlptracegrpc.Option{
				otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
			}
			if cfg.OTLPInsecure {
				opts = append(opts, otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
			}

			exporter, err = otlptracegrpc.New(ctx, opts...)
			if err != nil {
				return nil, fmt.Errorf("create otlp grpc exporter: %w", err)
			}
		default:
			return nil, fmt.Errorf("invalid otlp transport type: %s", cfg.OTLPTransportType)
		}
	default:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("create stdout exporter: %w", err)
		}
	}

	// Create resource
	res := resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.ServiceVersion),
		semconv.DeploymentEnvironment(cfg.Env),
	)

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set global TracerProvider
	otel.SetTracerProvider(tp)

	// Set global TextMapPropagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}
