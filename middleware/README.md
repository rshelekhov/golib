# Middleware

Protocol-agnostic middleware packages for gRPC and HTTP services.

## Packages

### Request ID (`middleware/requestid`)

Extract and propagate request IDs across services.

**Features:**

- Extracts request ID from gRPC metadata or HTTP headers
- Generates new request ID if not present (using ksuid)
- Stores request ID in context for use throughout request lifecycle
- Supports both unary and streaming gRPC interceptors
- HTTP middleware for standard HTTP handlers

### Logging (`middleware/logging`)

Log gRPC requests with method, status, and duration.

**Features:**

- Logs unary and streaming gRPC requests
- Includes method name, status code, and duration
- Configurable logger

### Recovery (`middleware/recovery`)

Recover from panics in gRPC handlers.

**Features:**

- Panic recovery for unary and streaming interceptors
- Logs panic details
- Returns internal server error instead of crashing

### Validation (`middleware/validation`)

Validate gRPC requests using protobuf validation.

**Features:**

- Automatic validation of requests implementing `Validate() error`
- Returns InvalidArgument status on validation failure

## Usage Example

```go
import (
    "github.com/rshelekhov/golib/middleware/logging"
    "github.com/rshelekhov/golib/middleware/recovery"
    "github.com/rshelekhov/golib/middleware/requestid"
    "github.com/rshelekhov/golib/middleware/validation"
)

// gRPC server setup
serverOpts := []grpc.ServerOption{
    grpc.ChainUnaryInterceptor(
        requestid.UnaryServerInterceptorFunc(),
        logging.UnaryServerInterceptor(logger),
        recovery.UnaryServerInterceptor(logger),
        validation.UnaryServerInterceptor(),
    ),
}

// HTTP middleware
mux := http.NewServeMux()
handler := requestid.HTTPMiddleware()(mux)
```

## Design Philosophy

These middleware packages are designed to be:

- **Protocol-agnostic**: Reusable across gRPC, HTTP, and future protocols
- **Composable**: Can be used independently or together
- **Reusable**: Importable in any Go project, not tied to specific server implementations
- **Minimal**: Each package has a single, focused responsibility
