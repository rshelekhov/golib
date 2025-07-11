package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// App represents the main application structure
type App struct {
	options     *Options
	grpcServer  *grpc.Server
	httpServer  *http.Server
	healthCheck *health.Server
	mux         *runtime.ServeMux
	httpMux     *http.ServeMux
}

// GRPCProvider is an interface for any service that can register with gRPC
type GRPCProvider interface {
	RegisterGRPC(*grpc.Server)
}

// HTTPProvider is an interface for services that can register HTTP handlers
type HTTPProvider interface {
	RegisterHTTP(context.Context, *runtime.ServeMux) error
}

// Service is a unified interface for services that can register with both gRPC and HTTP
type Service interface {
	GRPCProvider
	HTTPProvider
}

// ReadinessProvider is an interface for services that can report their readiness status
type ReadinessProvider interface {
	ReadinessChecks() []ReadinessCheck
}

// NewApp creates a new application instance with the given options
func NewApp(ctx context.Context, opts ...Option) (*App, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	if options.grpcPort <= 0 {
		return nil, fmt.Errorf("gRPC port must be specified and be greater than 0")
	}

	var httpServer *http.Server
	var httpMux *http.ServeMux
	var gwMux *runtime.ServeMux

	healthCheck := health.NewServer()

	serverOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(options.unaryInterceptors...),
		grpc.ChainStreamInterceptor(options.streamInterceptors...),
	}
	if options.statsHandler != nil {
		serverOpts = append(serverOpts, grpc.StatsHandler(options.statsHandler))
	}

	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(serverOpts...)

	// Register health check service
	healthpb.RegisterHealthServer(grpcServer, healthCheck)

	// Enable reflection for development tools
	if options.enableReflection {
		reflection.Register(grpcServer)
	}

	// Create HTTP server for gRPC-Gateway if port is specified
	if options.httpPort > 0 {
		// Create HTTP mux for gRPC-Gateway
		gwMux = runtime.NewServeMux(options.muxOptions...)

		// Create main HTTP mux for both gRPC-Gateway and other HTTP handlers
		httpMux = http.NewServeMux()

		// Register health check endpoints
		WithHealthEndpoints(httpMux, healthCheck)

		// Handle gRPC-Gateway requests
		httpMux.Handle("/", gwMux)

		// Create HTTP server with configured mux
		httpServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", options.httpPort),
			Handler: options.wrapHTTPHandler(httpMux),
		}
	}

	return &App{
		options:     options,
		grpcServer:  grpcServer,
		httpServer:  httpServer,
		healthCheck: healthCheck,
		mux:         gwMux,
		httpMux:     httpMux,
	}, nil
}

// Run starts the application servers and blocks until shutdown
func (a *App) Run(ctx context.Context, service GRPCProvider) error {
	// Set health check to serving
	a.healthCheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Create error group for concurrent server management
	g, ctx := errgroup.WithContext(ctx)

	// Start gRPC server
	a.startGRPCServer(g, service)

	// Register readiness checks if available
	if a.httpMux != nil && a.httpServer != nil {
		if readinessProvider, ok := service.(ReadinessProvider); ok {
			readinessChecks := readinessProvider.ReadinessChecks()
			a.httpMux.HandleFunc("/readyz", ReadinessHandler(readinessChecks...))
		}
	}

	// Start HTTP server if initialized
	if a.httpServer != nil && a.mux != nil {
		// Check if service implements HTTPProvider
		if httpProvider, ok := service.(HTTPProvider); ok {
			if err := a.startHTTPServer(ctx, g, httpProvider); err != nil {
				return err
			}
		} else {
			// Start HTTP server without registering HTTP handlers
			g.Go(func() error {
				a.options.logger.Info("starting HTTP server", "port", a.options.httpPort)
				if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
					return fmt.Errorf("HTTP server error: %w", err)
				}
				return nil
			})
		}
	}

	// Handle graceful shutdown
	a.handleGracefulShutdown(ctx, g)

	return g.Wait()
}

// startGRPCServer initializes and starts the gRPC server
func (a *App) startGRPCServer(g *errgroup.Group, service GRPCProvider) {
	// Register service with gRPC server
	service.RegisterGRPC(a.grpcServer)

	// Start gRPC server
	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.options.grpcPort))
		if err != nil {
			return fmt.Errorf("failed to listen on gRPC port: %w", err)
		}

		a.options.logger.Info("starting gRPC server", "port", a.options.grpcPort)
		if err := a.grpcServer.Serve(lis); err != nil {
			return fmt.Errorf("gRPC server error: %w", err)
		}
		return nil
	})
}

// startHTTPServer initializes and starts the HTTP server
func (a *App) startHTTPServer(ctx context.Context, g *errgroup.Group, provider HTTPProvider) error {
	// Register HTTP handlers
	if err := provider.RegisterHTTP(ctx, a.mux); err != nil {
		return fmt.Errorf("failed to register HTTP handlers: %w", err)
	}

	// Start HTTP server
	g.Go(func() error {
		a.options.logger.Info("starting HTTP server", "port", a.options.httpPort)
		if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("HTTP server error: %w", err)
		}
		return nil
	})

	return nil
}

// handleGracefulShutdown manages graceful shutdown on signals or context done
func (a *App) handleGracefulShutdown(ctx context.Context, g *errgroup.Group) {
	g.Go(func() error {
		// Create signal channel for shutdown
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case s := <-sigCh:
			a.options.logger.Info("received signal, shutting down", "signal", s.String())
		case <-ctx.Done():
			a.options.logger.Info("context done, shutting down")
		}

		a.Shutdown()
		return nil
	})
}

// Shutdown gracefully stops the application servers
func (a *App) Shutdown() {
	// Set health check to not serving
	a.healthCheck.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), a.options.shutdownTimeout)
	defer cancel()

	// Shutdown HTTP server
	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			a.options.logger.Error("error shutting down HTTP server", "error", err)
		}
	}

	// Gracefully stop gRPC server
	if a.grpcServer != nil {
		stopped := make(chan struct{})
		go func() {
			a.grpcServer.GracefulStop()
			close(stopped)
		}()

		// Force stop if graceful stop takes too long
		select {
		case <-ctx.Done():
			a.options.logger.Info("force stopping gRPC server")
			a.grpcServer.Stop()
		case <-stopped:
			a.options.logger.Info("gRPC server stopped gracefully")
		}
	}
}
