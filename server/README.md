# Server Bootstrap

Server bootstrap is part of the golib platform library for standardizing server initialization with gRPC and HTTP support. It provides a unified way to handle common server concerns:

- gRPC and HTTP server setup and configuration
- Kubernetes health and readiness probes
- Graceful shutdown
- Standard middleware and interceptors for logging, recovery, etc.
- Configuration via functional options

## Installation

```bash
go get github.com/rshelekhov/golib/server
```

## Basic Usage

```go
package main

import (
    "context"
    "log"
    "log/slog"
    "os"
    "time"

    "github.com/rshelekhov/golib/server"
    "google.golang.org/grpc"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// Implement the necessary interfaces for your service
type MyService struct{}

// RegisterGRPC registers your service with the gRPC server
func (s *MyService) RegisterGRPC(grpcServer *grpc.Server) {
    // Your gRPC service registration here
    // pb.RegisterYourServiceServer(grpcServer, &yourServiceImpl{})
}

// RegisterHTTP registers your service's HTTP handlers (optional)
func (s *MyService) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {
    // Your gRPC-Gateway registration here
    // return pb.RegisterYourServiceHandlerServer(ctx, mux, &yourServiceImpl{})
    return nil
}

func main() {
    ctx := context.Background()
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

    // Create the application with custom options
    app, err := server.NewApp(
        ctx,
        server.WithGRPCPort(9000),            // Required for gRPC server
        server.WithHTTPPort(8000),            // Optional, for HTTP Gateway
        server.WithLogger(logger),
        server.WithUnaryInterceptors(
            server.LoggingUnaryInterceptor(logger),
            server.RecoveryUnaryInterceptor(logger),
        ),
    )
    if err != nil {
        log.Fatalf("Failed to create application: %v", err)
    }

    // Start the servers with your service
    if err := app.Run(ctx, &MyService{}); err != nil {
        log.Fatalf("Application error: %v", err)
    }
}
```

## Configuration Options

Server bootstrap uses functional options for configuration:

- `WithGRPCPort(port int)` - Set the gRPC server port (required)
- `WithHTTPPort(port int)` - Set the HTTP server port (optional, for HTTP Gateway)
- `WithReflection(enable bool)` - Enable/disable gRPC reflection (default: enabled)
- `WithShutdownTimeout(timeout time.Duration)` - Set timeout for graceful shutdown (default: 10s)
- `WithUnaryInterceptors(...)` - Add gRPC unary interceptors
- `WithStreamInterceptors(...)` - Add gRPC stream interceptors
- `WithMuxOptions(...)` - Add gRPC-Gateway ServeMux options
- `WithHTTPMiddleware(...)` - Add HTTP middleware
- `WithLogger(logger *slog.Logger)` - Set the logger

## Server Modes

The library supports different server modes:

1. **gRPC only** - Only the gRPC server is started (default)
2. **gRPC with HTTP Gateway** - Both gRPC and HTTP Gateway servers are started

## Health Checks

The library automatically provides Kubernetes-compatible health endpoints:

- `/healthz` - Liveness probe to check if the service is running
- `/readyz` - Readiness probe to check if the service is ready to receive traffic

## Complete Example

See the `example/` directory for complete working examples.

## License

[MIT License](../LICENSE)
