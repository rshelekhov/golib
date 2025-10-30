package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rshelekhov/golib/middleware/cors"
	"github.com/rshelekhov/golib/middleware/logging"
	"github.com/rshelekhov/golib/middleware/recovery"
	"github.com/rshelekhov/golib/server"
	"google.golang.org/grpc"
)

// ExampleService is a simple implementation of the Service interface
type ExampleService struct{}

// RegisterGRPC registers the service with a gRPC server
func (s *ExampleService) RegisterGRPC(grpcServer *grpc.Server) {
	// In a real application, register your gRPC service here
	// For example: pb.RegisterYourServiceServer(grpcServer, &yourServiceImpl{})
	slog.Info("registered gRPC service")
}

// RegisterHTTP registers the service's HTTP handlers
func (s *ExampleService) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {
	// In a real application, register your HTTP handlers here
	// For example: if err := pb.RegisterYourServiceHandlerServer(ctx, mux, &yourServiceImpl{}); err != nil { return err }
	slog.Info("registered HTTP handlers")
	return nil
}

func main() {
	// Create a context
	ctx := context.Background()

	// Configure a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create the application with custom options
	app, err := server.NewApp(
		ctx,
		server.WithGRPCPort(9000),
		server.WithHTTPPort(8000),
		server.WithLogger(logger),
		server.WithShutdownTimeout(time.Second*5),
		server.WithUnaryInterceptors(
			logging.UnaryServerInterceptor(logger),
			recovery.UnaryServerInterceptor(logger),
		),
		server.WithStreamInterceptors(
			logging.StreamServerInterceptor(logger),
			recovery.StreamServerInterceptor(logger),
		),
		server.WithHTTPMiddleware(
			logging.Middleware(logger),
			recovery.Middleware(logger),
			cors.Middleware([]string{"*"}),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Create a service
	service := &ExampleService{}

	// Run the application
	logger.Info("starting application")
	if err := app.Run(ctx, service); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
