package server

import (
	"context"
	"net/http"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// HealthHandler creates an HTTP handler for health checks that uses the gRPC health check service
func HealthHandler(healthCheck *health.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the service is serving
		check, err := healthCheck.Check(context.Background(), &healthpb.HealthCheckRequest{})
		if err != nil || check.Status != healthpb.HealthCheckResponse_SERVING {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("Service Unavailable"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}

// ReadinessHandler creates an HTTP handler for readiness checks
func ReadinessHandler(checks ...ReadinessCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, check := range checks {
			if err := check.Check(r.Context()); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Ready"))
	}
}

// ReadinessCheck is an interface for components that can report their readiness status
type ReadinessCheck interface {
	Check(ctx context.Context) error
}

// WithHealthEndpoints adds standard health check endpoints to a ServeMux
func WithHealthEndpoints(mux *http.ServeMux, healthCheck *health.Server, readinessChecks ...ReadinessCheck) {
	// Live probe - is the service running?
	mux.HandleFunc("/healthz", HealthHandler(healthCheck))

	// Ready probe - is the service ready to receive traffic?
	mux.HandleFunc("/readyz", ReadinessHandler(readinessChecks...))
}
