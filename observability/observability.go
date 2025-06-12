package observability

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/rshelekhov/golib/observability/logger"
	"github.com/rshelekhov/golib/observability/metrics"
	"github.com/rshelekhov/golib/observability/tracing"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Env            string
	ServiceName    string
	ServiceVersion string
	EnableMetrics  bool
}

type Observability struct {
	Logger         *slog.Logger
	MetricsHandler http.Handler
	LoggerProvider *sdklog.LoggerProvider
	MeterProvider  *sdkmetric.MeterProvider
	TracerProvider *sdktrace.TracerProvider
}

// Init initializes observability with stdout tracing and Prometheus metrics
func Init(ctx context.Context, cfg Config) (*Observability, error) {
	// Initialize logger with stdout (for development)
	loggerProvider, logg, err := logger.InitLogger(cfg.ServiceName, cfg.ServiceVersion, cfg.Env)
	if err != nil {
		return nil, err
	}
	
	// Initialize tracing with stdout (for development)
	tracerProvider, err := tracing.InitTracer(cfg.ServiceName, cfg.ServiceVersion, cfg.Env)
	if err != nil {
		return nil, err
	}
	
	var metricsHandler http.Handler
	var meterProvider *sdkmetric.MeterProvider
	if cfg.EnableMetrics {
		// Default to Prometheus exporter
		meterProvider, metricsHandler, err = metrics.InitMeter(
			cfg.ServiceName,
			cfg.ServiceVersion,
			cfg.Env,
		)
		if err != nil {
			return nil, err
		}
	}
	
	return &Observability{
		Logger:         logg,
		MetricsHandler: metricsHandler,
		LoggerProvider: loggerProvider,
		MeterProvider:  meterProvider,
		TracerProvider: tracerProvider,
	}, nil
}

// InitWithOTLP initializes observability with OTLP exporters
func InitWithOTLP(ctx context.Context, cfg Config, otlpEndpoint string) (*Observability, error) {
	// Initialize logger with OTLP
	loggerProvider, logg, err := logger.InitLoggerOTLP(ctx, cfg.ServiceName, cfg.ServiceVersion, cfg.Env, otlpEndpoint)
	if err != nil {
		return nil, err
	}
	
	// Initialize tracing with OTLP
	tracerProvider, err := tracing.InitTracerOTLP(ctx, cfg.ServiceName, cfg.ServiceVersion, cfg.Env, otlpEndpoint)
	if err != nil {
		return nil, err
	}
	
	var meterProvider *sdkmetric.MeterProvider
	if cfg.EnableMetrics {
		meterProvider, err = metrics.InitMeterOTLP(
			ctx,
			cfg.ServiceName,
			cfg.ServiceVersion,
			cfg.Env,
			otlpEndpoint,
		)
		if err != nil {
			return nil, err
		}
	}
	
	return &Observability{
		Logger:         logg,
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
		return errs[0] // Return first error for simplicity
	}
	return nil
}

// HTTPMetricsMiddleware returns http.Handler with otel metrics
func HTTPMetricsMiddleware(next http.Handler) http.Handler { 
	return metrics.Middleware(next) 
}
