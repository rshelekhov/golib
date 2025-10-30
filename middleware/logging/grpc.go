package logging

import (
	"context"
	"log/slog"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor creates a gRPC unary interceptor for logging requests
func UnaryServerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		// Get method name
		method := path.Base(info.FullMethod)

		// Call the handler
		resp, err = handler(ctx, req)

		// Get status code
		statusCode := codes.OK
		if err != nil {
			if s, ok := status.FromError(err); ok {
				statusCode = s.Code()
			}
		}

		// Log request
		logger.Info("grpc request",
			"method", method,
			"status", statusCode.String(),
			"duration", time.Since(start),
		)

		return resp, err
	}
}

// StreamServerInterceptor creates a gRPC stream interceptor for logging requests
func StreamServerInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Get method name
		method := path.Base(info.FullMethod)

		// Call the handler
		err := handler(srv, ss)

		// Get status code
		statusCode := codes.OK
		if err != nil {
			if s, ok := status.FromError(err); ok {
				statusCode = s.Code()
			}
		}

		// Log request
		logger.Info("grpc stream",
			"method", method,
			"status", statusCode.String(),
			"duration", time.Since(start),
		)

		return err
	}
}
