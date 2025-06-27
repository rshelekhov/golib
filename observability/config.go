package observability

import (
	"fmt"
	"log/slog"
	"strings"
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

type Config struct {
	Env            string
	ServiceName    string
	ServiceVersion string
	EnableMetrics  bool
	OTLPEndpoint   string
	LogLevel       slog.Level
}

func (c Config) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}
	if c.ServiceVersion == "" {
		return fmt.Errorf("service version is required")
	}
	if c.Env == "" {
		return fmt.Errorf("environment is required")
	}
	if _, ok := supportedEnvs[c.Env]; !ok {
		return fmt.Errorf("unsupported environment: %s (supported: %s)", c.Env, strings.Join(getSupportedEnvs(), ", "))
	}
	if c.requiresOTLPEndpoint() && c.OTLPEndpoint == "" {
		return fmt.Errorf("OTLP endpoint is required for environment %s", c.Env)
	}
	return nil
}

func (c Config) requiresOTLPEndpoint() bool {
	return c.Env == EnvDev || c.Env == EnvProd
}

func getSupportedEnvs() []string {
	envs := make([]string, 0, len(supportedEnvs))
	for env := range supportedEnvs {
		envs = append(envs, env)
	}
	return envs
}

func getDefaultLogLevel(env string) slog.Level {
	switch env {
	case EnvLocal:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

// NewConfig creates config with environment-based defaults and optional log level override
func NewConfig(env, serviceName, serviceVersion string, enableMetrics bool, otlpEndpoint string, logLevel ...slog.Level) (Config, error) {
	level := getDefaultLogLevel(env)
	if len(logLevel) > 0 {
		level = logLevel[0]
	}

	cfg := Config{
		Env:            env,
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		EnableMetrics:  enableMetrics,
		OTLPEndpoint:   otlpEndpoint,
		LogLevel:       level,
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
