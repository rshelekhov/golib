package server

import (
	"context"
	"log/slog"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingUnaryInterceptor creates a gRPC unary interceptor for logging requests
func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
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

// RecoveryUnaryInterceptor creates a gRPC unary interceptor for recovering from panics
func RecoveryUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
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

// ValidationUnaryInterceptor creates a gRPC unary interceptor for validating requests
func ValidationUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if validator, ok := req.(interface{ Validate() error }); ok {
			if err := validator.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		return handler(ctx, req)
	}
}

// LoggingStreamInterceptor creates a gRPC stream interceptor for logging requests
func LoggingStreamInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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

// RecoveryStreamInterceptor creates a gRPC stream interceptor for recovering from panics
func RecoveryStreamInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
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
