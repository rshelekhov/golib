package observability

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/rshelekhov/golib/observability/tracing"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

var supportedEnvs = map[string]struct{}{
	EnvLocal: {},
	EnvDev:   {},
	EnvProd:  {},
}

var supportedOTLPTransportTypes = map[tracing.OTLPTransportType]struct{}{
	tracing.OTLPTransportGRPC: {},
	tracing.OTLPTransportHTTP: {},
}

type Config struct {
	Env               string
	ServiceName       string
	ServiceVersion    string
	EnableMetrics     bool
	OTLPEndpoint      string
	OTLPTransportType tracing.OTLPTransportType
	LogLevel          slog.Level

	// TLS configuration for OTLP exporters
	// If true, uses TLS (default for production)
	// If false, uses insecure connection (useful for local development)
	OTLPInsecure bool
}

type ConfigParams struct {
	Env               string
	ServiceName       string
	ServiceVersion    string
	EnableMetrics     bool
	OTLPEndpoint      string
	OTLPTransportType string
	OTLPInsecure      *bool // Use pointer to distinguish between "not set" and "explicitly false"
}

func (c ConfigParams) Validate() error {
	var errMessages []string

	if c.ServiceName == "" {
		errMessages = append(errMessages, "service name is required")
	}
	if c.ServiceVersion == "" {
		errMessages = append(errMessages, "service version is required")
	}

	if c.Env == "" {
		errMessages = append(errMessages, "environment is required")
	}
	if _, ok := supportedEnvs[c.Env]; !ok {
		errMessages = append(errMessages, fmt.Sprintf("unsupported environment: %s (supported: %s)", c.Env, strings.Join(getSupportedEnvs(), ", ")))
	}
	if c.requiresOTLPEndpoint() && c.OTLPEndpoint == "" {
		errMessages = append(errMessages, fmt.Sprintf("OTLP endpoint is required for environment %s", c.Env))
	}
	if c.requiresOTLPEndpoint() && c.OTLPTransportType == "" {
		errMessages = append(errMessages, fmt.Sprintf("OTLP transport type is required for environment %s", c.Env))
	}
	if c.OTLPTransportType != "" && !isValidOTLPTransportType(c.OTLPTransportType) {
		errMessages = append(errMessages, fmt.Sprintf("unsupported OTLP transport type: %s (supported: %s)", c.OTLPTransportType, strings.Join(getSupportedOTLPTransportTypes(), ", ")))
	}

	if len(errMessages) > 0 {
		return fmt.Errorf("%s", strings.Join(errMessages, "; "))
	}
	return nil
}

func isValidOTLPTransportType(transportType string) bool {
	normalizedType := tracing.OTLPTransportType(strings.ToLower(transportType))
	_, ok := supportedOTLPTransportTypes[normalizedType]
	return ok
}

func (c ConfigParams) requiresOTLPEndpoint() bool {
	return c.Env == EnvDev || c.Env == EnvProd
}

func getSupportedEnvs() []string {
	envs := make([]string, 0, len(supportedEnvs))
	for env := range supportedEnvs {
		envs = append(envs, env)
	}
	return envs
}

func getSupportedOTLPTransportTypes() []string {
	types := make([]string, 0, len(supportedOTLPTransportTypes))
	for t := range supportedOTLPTransportTypes {
		types = append(types, string(t))
	}
	return types
}

func getDefaultLogLevel(env string) slog.Level {
	switch env {
	case EnvLocal:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

func getDefaultOTLPInsecure(env string) bool {
	switch env {
	case EnvLocal, EnvDev:
		return true // Use insecure connections for local/dev
	default:
		return false // Use TLS for production
	}
}

// Option defines a functional option for Config
type Option func(*Config)

// WithLogLevel sets a custom log level, overriding environment defaults
func WithLogLevel(level slog.Level) Option {
	return func(cfg *Config) {
		cfg.LogLevel = level
	}
}

// WithOTLPInsecure sets whether to use insecure OTLP connections
func WithOTLPInsecure(insecure bool) Option {
	return func(cfg *Config) {
		cfg.OTLPInsecure = insecure
	}
}

// NewConfig creates config with environment-based defaults and optional overrides
func NewConfig(params ConfigParams, opts ...Option) (Config, error) {
	if err := params.Validate(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		Env:               params.Env,
		ServiceName:       params.ServiceName,
		ServiceVersion:    params.ServiceVersion,
		EnableMetrics:     params.EnableMetrics,
		OTLPEndpoint:      params.OTLPEndpoint,
		OTLPTransportType: tracing.OTLPTransportType(params.OTLPTransportType),
		LogLevel:          getDefaultLogLevel(params.Env),
		OTLPInsecure:      getDefaultOTLPInsecure(params.Env),
	}

	// If user explicitly set OTLPInsecure in params, use that instead of default
	if params.OTLPInsecure != nil {
		cfg.OTLPInsecure = *params.OTLPInsecure
	}

	// Apply options
	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg, nil
}
