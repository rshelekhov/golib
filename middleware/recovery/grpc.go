package recovery

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor creates a gRPC unary interceptor for recovering from panics
func UnaryServerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("grpc server panic recovered",
					"error", r,
					"method", info.FullMethod,
				)

				err = status.Error(codes.Internal, "Internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// StreamServerInterceptor creates a gRPC stream interceptor for recovering from panics
func StreamServerInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("grpc stream server panic recovered",
					"error", r,
					"method", info.FullMethod,
				)

				_ = status.Error(codes.Internal, "Internal server error")
			}
		}()

		return handler(srv, ss)
	}
}
