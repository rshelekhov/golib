package observability

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/rshelekhov/golib/observability/logger"
	"github.com/rshelekhov/golib/observability/metrics"
	"github.com/rshelekhov/golib/observability/tracing"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Observability struct {
	Logger         *slog.Logger
	MetricsHandler http.Handler
	LoggerProvider *sdklog.LoggerProvider
	MeterProvider  *sdkmetric.MeterProvider
	TracerProvider *sdktrace.TracerProvider
}

// Init initializes observability with automatic exporter selection
func Init(ctx context.Context, cfg Config) (*Observability, error) {
	// Determine if we should use OTLP based on configuration
	useOTLP := cfg.OTLPEndpoint != "" && cfg.Env != EnvLocal

	// Initialize logger
	loggerCfg := logger.Config{
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Env:            cfg.Env,
		Level:          cfg.LogLevel,
	}
	if useOTLP {
		loggerCfg.Endpoint = cfg.OTLPEndpoint
	}
	loggerProvider, log, err := logger.Init(ctx, loggerCfg)
	if err != nil {
		return nil, err
	}

	// Initialize tracing
	tracingCfg := tracing.Config{
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Env:            cfg.Env,
	}
	if useOTLP {
		tracingCfg.ExporterType = tracing.ExporterOTLP
		tracingCfg.OTLPEndpoint = cfg.OTLPEndpoint
		tracingCfg.OTLPTransportType = cfg.OTLPTransportType
	} else {
		tracingCfg.ExporterType = tracing.ExporterStdout
	}
	tracerProvider, err := tracing.Init(ctx, tracingCfg)
	if err != nil {
		return nil, err
	}

	var metricsHandler http.Handler
	var meterProvider *sdkmetric.MeterProvider

	// Metrics are completely disabled for local development
	// For other environments, respect the EnableMetrics flag
	if cfg.Env != EnvLocal && cfg.EnableMetrics {
		metricsCfg := metrics.Config{
			ServiceName:    cfg.ServiceName,
			ServiceVersion: cfg.ServiceVersion,
			Env:            cfg.Env,
		}
		if useOTLP {
			metricsCfg.ExporterType = metrics.ExporterOTLP
			metricsCfg.OTLPEndpoint = cfg.OTLPEndpoint
		} else {
			metricsCfg.ExporterType = metrics.ExporterPrometheus
		}
		meterProvider, metricsHandler, err = metrics.Init(ctx, metricsCfg)
		if err != nil {
			return nil, err
		}
	}

	return &Observability{
		Logger:         log,
		MetricsHandler: metricsHandler,
		LoggerProvider: loggerProvider,
		MeterProvider:  meterProvider,
		TracerProvider: tracerProvider,
	}, nil
}

// Shutdown gracefully shuts down all observability components
func (o *Observability) Shutdown(ctx context.Context) error {
	var errs []error

	if o.TracerProvider != nil {
		if err := o.TracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if o.MeterProvider != nil {
		if err := o.MeterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if o.LoggerProvider != nil {
		if err := o.LoggerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// HTTPMetricsMiddleware returns http.Handler with otel metrics
func HTTPMetricsMiddleware(next http.Handler) http.Handler {
	return metrics.Middleware(next)
}
