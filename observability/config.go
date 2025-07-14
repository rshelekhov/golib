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
	tracing.OTLPGRPC: {},
	tracing.OTLPHTTP: {},
}

type Config struct {
	Env               string
	ServiceName       string
	ServiceVersion    string
	EnableMetrics     bool
	OTLPEndpoint      string
	OTLPTransportType tracing.OTLPTransportType
	LogLevel          slog.Level
}

type ConfigParams struct {
	Env               string
	ServiceName       string
	ServiceVersion    string
	EnableMetrics     bool
	OTLPEndpoint      string
	OTLPTransportType tracing.OTLPTransportType
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
		errMessages = append(errMessages, "OTLP transport type is required for environment %s", c.Env)
	}
	if _, ok := supportedOTLPTransportTypes[c.OTLPTransportType]; !ok {
		errMessages = append(errMessages, fmt.Sprintf("unsupported OTLP transport type: %s (supported: %s)", c.OTLPTransportType, strings.Join(getSupportedOTLPTransportTypes(), ", ")))
		return fmt.Errorf("%s", strings.Join(errMessages, "; "))
	}

	if len(errMessages) > 0 {
		return fmt.Errorf("%s", strings.Join(errMessages, "; "))
	}
	return nil
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

// Option defines a functional option for Config
type Option func(*Config)

// WithLogLevel sets a custom log level, overriding environment defaults
func WithLogLevel(level slog.Level) Option {
	return func(cfg *Config) {
		cfg.LogLevel = level
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
		OTLPTransportType: params.OTLPTransportType,
		LogLevel:          getDefaultLogLevel(params.Env),
	}

	// Apply options
	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg, nil
}
