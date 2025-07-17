package observability

import (
	"testing"

	"github.com/rshelekhov/golib/observability/tracing"
)

func TestTLSConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		env           string
		expectedTLS   bool
		overrideTLS   *bool
		finalExpected bool
	}{
		{
			name:          "Local environment defaults to insecure",
			env:           EnvLocal,
			expectedTLS:   true, // insecure = true means no TLS
			overrideTLS:   nil,
			finalExpected: true,
		},
		{
			name:          "Dev environment defaults to insecure",
			env:           EnvDev,
			expectedTLS:   true, // insecure = true means no TLS
			overrideTLS:   nil,
			finalExpected: true,
		},
		{
			name:          "Prod environment defaults to secure",
			env:           EnvProd,
			expectedTLS:   false, // insecure = false means TLS enabled
			overrideTLS:   nil,
			finalExpected: false,
		},
		{
			name:          "Override dev to use TLS",
			env:           EnvDev,
			expectedTLS:   true,              // default for dev
			overrideTLS:   &[]bool{false}[0], // override to use TLS
			finalExpected: false,
		},
		{
			name:          "Override prod to disable TLS",
			env:           EnvProd,
			expectedTLS:   false,            // default for prod
			overrideTLS:   &[]bool{true}[0], // override to disable TLS
			finalExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ConfigParams{
				Env:               tt.env,
				ServiceName:       "test-service",
				ServiceVersion:    "1.0.0",
				EnableMetrics:     true,
				OTLPEndpoint:      "localhost:4317",
				OTLPTransportType: tracing.OTLPGRPC,
			}

			var opts []Option
			if tt.overrideTLS != nil {
				opts = append(opts, WithOTLPInsecure(*tt.overrideTLS))
			}

			cfg, err := NewConfig(params, opts...)
			if err != nil {
				t.Fatalf("NewConfig failed: %v", err)
			}

			if cfg.OTLPInsecure != tt.finalExpected {
				t.Errorf("Expected OTLPInsecure=%v, got %v", tt.finalExpected, cfg.OTLPInsecure)
			}
		})
	}
}

func TestTLSConfigurationDefaults(t *testing.T) {
	// Test that getDefaultOTLPInsecure returns correct values
	tests := []struct {
		env      string
		expected bool
	}{
		{EnvLocal, true}, // insecure for local
		{EnvDev, true},   // insecure for dev
		{EnvProd, false}, // secure for prod
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			result := getDefaultOTLPInsecure(tt.env)
			if result != tt.expected {
				t.Errorf("getDefaultOTLPInsecure(%s) = %v, expected %v", tt.env, result, tt.expected)
			}
		})
	}
}

func TestConfigParamsWithExplicitTLS(t *testing.T) {
	// Test that explicitly setting OTLPInsecure in ConfigParams works
	insecureTrue := true
	params := ConfigParams{
		Env:               EnvProd, // defaults to secure
		ServiceName:       "test-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "localhost:4317",
		OTLPTransportType: tracing.OTLPGRPC,
		OTLPInsecure:      &insecureTrue, // explicitly set to insecure
	}

	cfg, err := NewConfig(params)
	if err != nil {
		t.Fatalf("NewConfig failed: %v", err)
	}

	if !cfg.OTLPInsecure {
		t.Errorf("Expected OTLPInsecure=true (explicit override), got %v", cfg.OTLPInsecure)
	}
}

func TestFunctionalOptionOverride(t *testing.T) {
	// Test that functional options override both defaults and explicit params
	insecureFalse := false
	params := ConfigParams{
		Env:               EnvProd,
		ServiceName:       "test-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "localhost:4317",
		OTLPTransportType: tracing.OTLPGRPC,
		OTLPInsecure:      &insecureFalse, // explicitly secure
	}

	// Override with functional option
	cfg, err := NewConfig(params, WithOTLPInsecure(true))
	if err != nil {
		t.Fatalf("NewConfig failed: %v", err)
	}

	if !cfg.OTLPInsecure {
		t.Errorf("Expected functional option to override: OTLPInsecure=true, got %v", cfg.OTLPInsecure)
	}
}
