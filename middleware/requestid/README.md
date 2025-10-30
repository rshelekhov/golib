# Request ID Middleware

Middleware for extracting and propagating request IDs across gRPC and HTTP services.

## Features

- Extracts request ID from incoming gRPC metadata or HTTP headers
- Generates new request ID if not present (using ksuid)
- Stores request ID in context for use throughout request lifecycle
- Provides utilities for extracting request ID from context
- Supports both unary and streaming gRPC interceptors
- HTTP middleware for standard HTTP handlers

## Usage

### gRPC Unary Interceptor

```go
import "github.com/rshelekhov/golib/middleware/requestid"

// Using the Interceptor struct
interceptor := requestid.NewInterceptor()
serverOpts := []grpc.ServerOption{
    grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()),
}

// Or using the convenience function
serverOpts := []grpc.ServerOption{
    grpc.UnaryInterceptor(requestid.UnaryServerInterceptorFunc()),
}
```

### gRPC Stream Interceptor

```go
import "github.com/rshelekhov/golib/middleware/requestid"

interceptor := requestid.NewInterceptor()
serverOpts := []grpc.ServerOption{
    grpc.StreamInterceptor(interceptor.StreamServerInterceptor()),
}
```

### HTTP Middleware

```go
import "github.com/rshelekhov/golib/middleware/requestid"

mux := http.NewServeMux()
handler := requestid.HTTPMiddleware()(mux)
```

### Extracting Request ID

```go
import "github.com/rshelekhov/golib/middleware/requestid"

func myHandler(ctx context.Context) {
    requestID, ok := requestid.FromContext(ctx)
    if ok {
        // Use requestID
    }
}
```

## Constants

- `Header`: The header name used for request ID (`X-Request-ID`)
- `CtxKey`: The context key used to store request ID (`RequestID`)
