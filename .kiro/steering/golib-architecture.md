# GoLib Architecture and Standards

This is a Go library collection (golib) that provides standardized components for building microservices and applications. The repository follows a modular architecture with independent Go modules.

## Repository Structure

The repository uses Go workspaces with the following modules:

- `config/` - Configuration loading with auto-discovery
- `db/mongo/` - MongoDB wrapper with transaction management
- `db/postgres/pgxv5/` - PostgreSQL wrapper with pgx v5
- `db/redis/` - Redis wrapper with pipeline/transaction support
- `db/s3/` - S3/MinIO wrapper for object storage
- `observability/` - Logging, tracing, and metrics (OpenTelemetry)
- `server/` - gRPC/HTTP server bootstrap

## Go Module Standards

### Module Structure

Each module follows this pattern:

```
module-name/
├── go.mod                 # Independent module
├── go.sum
├── README.md             # Comprehensive documentation
├── CHANGELOG.md          # Version history
├── main.go or *.go       # Core implementation
├── *_test.go            # Unit tests
├── examples/            # Usage examples
├── testutil/            # Testing utilities (if needed)
└── infrastructure/      # Docker/deployment files (if needed)
```

### Naming Conventions

- Package names: lowercase, single word when possible
- Interface names: descriptive, often ending with -er (e.g., `Connection`, `TransactionManager`)
- Constructor functions: `New*` (e.g., `NewConnection`, `NewTransactionManager`)
- Configuration functions: `With*` for functional options (e.g., `WithTimeout`, `WithTracing`)

### Configuration Patterns

All modules use functional options pattern:

```go
// ConfigParams struct for required parameters
type ConfigParams struct {
    RequiredField string
    AnotherField  int
}

// Functional options for optional parameters
func WithOptionalSetting(value string) Option {
    return func(cfg *Config) {
        cfg.OptionalSetting = value
    }
}

// Constructor with validation
func NewConnection(ctx context.Context, params ConfigParams, opts ...Option) (*Connection, error) {
    if err := params.Validate(); err != nil {
        return nil, err
    }
    // Apply options and create connection
}
```

## Database Layer Standards

### Connection Management

- All database modules provide `NewConnection()` constructor
- Support context for cancellation and timeouts
- Implement proper connection pooling where applicable
- Provide `Close()` method for cleanup

### Transaction Management

- Implement `TransactionManager` interface
- Support nested transactions where possible
- Provide `RunTransaction()` methods with different isolation levels
- Use context to pass transaction state

### OpenTelemetry Integration

- Enable tracing by default with `WithTracing(true)`
- Allow disabling tracing with `WithTracing(false)`
- Instrument all database operations
- Follow OpenTelemetry semantic conventions

## Observability Standards

### Logging

- Use structured logging with `slog`
- Environment-based configuration (local=debug, prod=info)
- Pretty colorized output for local development
- Automatic trace_id/span_id injection

### Tracing

- OpenTelemetry integration by default
- Support both stdout and OTLP exporters
- Automatic exporter selection based on environment
- Proper span naming and attributes

### Metrics

- OpenTelemetry metrics with Prometheus export
- Disabled by default in local environment
- OTLP export for production environments
- Standard metric naming conventions

### TLS Configuration

- Smart environment-based defaults (insecure for local/dev, secure for prod)
- Configurable via functional options
- Support both gRPC and HTTP transports

## Server Standards

### gRPC/HTTP Bootstrap

- Unified server initialization
- Support for both gRPC-only and gRPC+HTTP Gateway modes
- Built-in health checks (`/healthz`, `/readyz`)
- Graceful shutdown with configurable timeout
- Standard middleware and interceptors

### Configuration

- Functional options for all server settings
- Required gRPC port, optional HTTP port
- Configurable interceptors and middleware
- OpenTelemetry integration support
