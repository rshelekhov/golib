package main

import (
	"context"
	"log"
	"net/http"

	"github.com/rshelekhov/golib/observability"
	"github.com/rshelekhov/golib/observability/metrics"
	"github.com/rshelekhov/golib/observability/tracing"
)

func main() {
	log.Println("=== TLS Configuration Examples ===")

	// Example 1: Production with TLS (default)
	prodWithTLSExample()

	// Example 2: Production with insecure connection (override)
	prodInsecureExample()

	// Example 3: Local development (insecure by default)
	localExample()

	// Example 4: Using functional options to override defaults
	functionalOptionsExample()
}

func prodWithTLSExample() {
	log.Println("\n1. Production with TLS (default behavior)")

	cfg, err := observability.NewConfig(observability.ConfigParams{
		Env:               observability.EnvProd,
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "otel-collector.company.com:4317",
		OTLPTransportType: tracing.OTLPGRPC,
		// OTLPInsecure is false by default for production
	})
	if err != nil {
		log.Printf("Config error: %v", err)
		return
	}

	log.Printf("Production config - OTLPInsecure: %v (uses TLS)", cfg.OTLPInsecure)

	// This would connect to OTLP collector using TLS
	// obs, err := observability.Init(context.Background(), cfg)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// defer obs.Shutdown(context.Background())
}

func prodInsecureExample() {
	log.Println("\n2. Production with insecure connection (explicit override)")

	cfg, err := observability.NewConfig(observability.ConfigParams{
		Env:               observability.EnvProd,
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "localhost:4317", // Local OTLP collector
		OTLPTransportType: tracing.OTLPGRPC,
		OTLPInsecure:      &[]bool{true}[0], // Explicitly set to insecure
	})
	if err != nil {
		log.Printf("Config error: %v", err)
		return
	}

	log.Printf("Production config - OTLPInsecure: %v (overridden to insecure)", cfg.OTLPInsecure)
}

func localExample() {
	log.Println("\n3. Local development (insecure by default)")

	cfg, err := observability.NewConfig(observability.ConfigParams{
		Env:               observability.EnvLocal,
		ServiceName:       "my-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     false,            // Metrics disabled for local
		OTLPEndpoint:      "",               // No OTLP for local
		OTLPTransportType: tracing.OTLPGRPC, // Required for validation, but not used
		// OTLPInsecure is true by default for local
	})
	if err != nil {
		log.Printf("Config error: %v", err)
		return
	}

	log.Printf("Local config - OTLPInsecure: %v (default for local)", cfg.OTLPInsecure)
}

func functionalOptionsExample() {
	log.Println("\n4. Using functional options to override TLS settings")

	// Start with dev environment (insecure by default)
	// But override to use TLS
	cfg, err := observability.NewConfig(
		observability.ConfigParams{
			Env:               observability.EnvDev,
			ServiceName:       "my-service",
			ServiceVersion:    "1.0.0",
			EnableMetrics:     true,
			OTLPEndpoint:      "secure-otel-collector.dev.company.com:4317",
			OTLPTransportType: tracing.OTLPGRPC,
		},
		observability.WithOTLPInsecure(false), // Override to use TLS
	)
	if err != nil {
		log.Printf("Config error: %v", err)
		return
	}

	log.Printf("Dev config with TLS override - OTLPInsecure: %v", cfg.OTLPInsecure)

	// Example of the opposite: force insecure for production (not recommended)
	cfg2, err := observability.NewConfig(
		observability.ConfigParams{
			Env:               observability.EnvProd,
			ServiceName:       "my-service",
			ServiceVersion:    "1.0.0",
			EnableMetrics:     true,
			OTLPEndpoint:      "localhost:4317",
			OTLPTransportType: tracing.OTLPGRPC,
		},
		observability.WithOTLPInsecure(true), // Force insecure (not recommended for prod)
	)
	if err != nil {
		log.Printf("Config error: %v", err)
		return
	}

	log.Printf("Prod config with insecure override - OTLPInsecure: %v (not recommended)", cfg2.OTLPInsecure)
}

// Example of a complete working application with TLS configuration
func completeExample() {
	cfg, err := observability.NewConfig(observability.ConfigParams{
		Env:               observability.EnvDev,
		ServiceName:       "tls-demo-service",
		ServiceVersion:    "1.0.0",
		EnableMetrics:     true,
		OTLPEndpoint:      "localhost:4317",
		OTLPTransportType: tracing.OTLPGRPC,
		// Will use insecure connection by default for dev environment
	})
	if err != nil {
		log.Fatal(err)
	}

	obs, err := observability.Init(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := obs.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	// HTTP server with observability
	http.Handle("/", metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obs.Logger.InfoContext(r.Context(), "handling request",
			"path", r.URL.Path,
			"tls_config", cfg.OTLPInsecure,
		)

		if _, err := w.Write([]byte("Hello with configurable TLS!")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})))

	log.Printf("Server starting with OTLP insecure: %v", cfg.OTLPInsecure)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
