# golib

A collection of reusable Go libraries for building production-ready services.

## Packages

### [config](config/)

Flexible configuration loader with auto-discovery, YAML and .env support, and environment variable integration.

### [server](server/)

Server bootstrap library for gRPC and HTTP services with health checks, graceful shutdown, and standard middleware.

### [middleware](middleware/)

Protocol-agnostic middleware packages:

- **requestid** - Request ID extraction and propagation
- **logging** - Request logging for gRPC and HTTP
- **recovery** - Panic recovery middleware
- **validation** - Request validation
- **cors** - CORS handling for HTTP

### [observability](observability/)

Observability tools for modern Go services:

- **logger** - Structured logging with slog
- **metrics** - Prometheus metrics for gRPC and HTTP
- **tracing** - OpenTelemetry tracing support

### [db](db/)

Database connection and transaction management:

- **mongo** - MongoDB client with transaction support
- **postgres/pgxv5** - PostgreSQL client using pgx v5
- **redis** - Redis client
- **s3** - AWS S3 client

## Installation

Install individual packages as needed:

```bash
go get github.com/rshelekhov/golib/server
go get github.com/rshelekhov/golib/middleware/requestid
go get github.com/rshelekhov/golib/config
# ... etc
```

## License

See [LICENSE](LICENSE) file for details.
