package requestid

import (
	"context"

	"github.com/segmentio/ksuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Interceptor handles request ID extraction and injection for gRPC
type Interceptor struct{}

// NewInterceptor creates a new request ID interceptor
func NewInterceptor() *Interceptor {
	return &Interceptor{}
}

// UnaryServerInterceptor returns a gRPC unary server interceptor that extracts
// or generates a request ID and adds it to the context
func (i *Interceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		requestID := extractFromGRPC(ctx)
		ctx = WithContext(ctx, requestID)

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that extracts
// or generates a request ID and adds it to the context
func (i *Interceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		requestID := extractFromGRPC(ctx)
		ctx = WithContext(ctx, requestID)

		// Wrap the server stream to carry the new context
		wrapped := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		return handler(srv, wrapped)
	}
}

// UnaryServerInterceptorFunc returns a gRPC unary server interceptor function
// for convenience when you don't need the Interceptor struct
func UnaryServerInterceptorFunc() grpc.UnaryServerInterceptor {
	return NewInterceptor().UnaryServerInterceptor()
}

// StreamServerInterceptorFunc returns a gRPC stream server interceptor function
// for convenience when you don't need the Interceptor struct
func StreamServerInterceptorFunc() grpc.StreamServerInterceptor {
	return NewInterceptor().StreamServerInterceptor()
}

// extractFromGRPC extracts the request ID from gRPC metadata or generates a new one
func extractFromGRPC(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newID()
	}

	values := md.Get(Header)
	if len(values) == 0 {
		return newID()
	}

	return values[0]
}

// newID generates a new request ID using ksuid
func newID() string {
	return ksuid.New().String()
}

// wrappedServerStream wraps grpc.ServerStream to override the context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
