package server

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

// Options holds all configuration for the application
type Options struct {
	// Server configuration
	grpcPort         int
	httpPort         int
	enableReflection bool
	shutdownTimeout  time.Duration

	// Middleware and interceptors
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	muxOptions         []runtime.ServeMuxOption
	httpMiddleware     []func(http.Handler) http.Handler

	// Tracing
	statsHandler stats.Handler

	// Logger
	logger *slog.Logger
}

// Option is a function that modifies Options
type Option func(*Options)

// defaultOptions returns the default configuration
func defaultOptions() *Options {
	return &Options{
		grpcPort:           9000,
		httpPort:           0,
		enableReflection:   true,
		shutdownTimeout:    time.Second * 10,
		unaryInterceptors:  []grpc.UnaryServerInterceptor{},
		streamInterceptors: []grpc.StreamServerInterceptor{},
		muxOptions:         []runtime.ServeMuxOption{},
		httpMiddleware:     []func(http.Handler) http.Handler{},
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}
}

// WithGRPCPort sets the gRPC server port
func WithGRPCPort(port int) Option {
	return func(o *Options) {
		o.grpcPort = port
	}
}

// WithHTTPPort sets the HTTP server port
func WithHTTPPort(port int) Option {
	return func(o *Options) {
		o.httpPort = port
	}
}

// WithReflection enables/disables gRPC reflection
func WithReflection(enable bool) Option {
	return func(o *Options) {
		o.enableReflection = enable
	}
}

// WithShutdownTimeout sets the timeout for graceful shutdown
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.shutdownTimeout = timeout
	}
}

// WithUnaryInterceptors adds gRPC unary interceptors
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(o *Options) {
		o.unaryInterceptors = append(o.unaryInterceptors, interceptors...)
	}
}

// WithStreamInterceptors adds gRPC stream interceptors
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(o *Options) {
		o.streamInterceptors = append(o.streamInterceptors, interceptors...)
	}
}

// WithMuxOptions adds gRPC-Gateway ServeMux options
func WithMuxOptions(options ...runtime.ServeMuxOption) Option {
	return func(o *Options) {
		o.muxOptions = append(o.muxOptions, options...)
	}
}

// WithHTTPMiddleware adds HTTP middleware
func WithHTTPMiddleware(middleware ...func(http.Handler) http.Handler) Option {
	return func(o *Options) {
		o.httpMiddleware = append(o.httpMiddleware, middleware...)
	}
}

// WithLogger sets the logger
func WithLogger(logger *slog.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

// WithStatsHandler sets the stats handler
func WithStatsHandler(statsHandler stats.Handler) Option {
	return func(o *Options) {
		o.statsHandler = statsHandler
	}
}

// wrapHTTPHandler applies all registered HTTP middleware to the handler
func (o *Options) wrapHTTPHandler(handler http.Handler) http.Handler {
	// Apply middleware in reverse order (last added is outermost)
	h := handler
	for i := len(o.httpMiddleware) - 1; i >= 0; i-- {
		h = o.httpMiddleware[i](h)
	}
	return h
}
