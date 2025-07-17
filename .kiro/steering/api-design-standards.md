# API Design and Versioning Standards

## Public API Design Principles

### Interface Design

- Design interfaces first, implementations second
- Keep interfaces small and focused (Interface Segregation Principle)
- Use composition over inheritance
- Prefer accepting interfaces, returning concrete types
- Make interfaces discoverable through embedding

```go
// Good: Small, focused interface
type Reader interface {
    Read(ctx context.Context, key string) ([]byte, error)
}

type Writer interface {
    Write(ctx context.Context, key string, data []byte) error
}

// Compose interfaces when needed
type ReadWriter interface {
    Reader
    Writer
}
```

### Function Signatures

- Use context.Context as the first parameter for operations that can be cancelled
- Return errors as the last return value
- Use functional options for optional parameters
- Prefer explicit parameters over configuration structs for simple cases

```go
// Good: Clear, extensible API
func NewConnection(ctx context.Context, dsn string, opts ...Option) (*Connection, error)

// Good: Context first, error last
func (c *Connection) Query(ctx context.Context, query string, args ...interface{}) (*Result, error)
```

### Configuration API

- Follow the functional options pattern defined in golib-architecture.md
- Provide sensible defaults for all optional parameters
- Validate configuration early and return clear errors
- Use ConfigParams struct for required parameters, functional options for optional ones

## Backward Compatibility

### Versioning Strategy

- Use semantic versioning (semver) for all modules
- Major version changes for breaking changes
- Minor version changes for new features
- Patch version changes for bug fixes
- Document breaking changes clearly in CHANGELOG.md

### Breaking Changes

- Avoid breaking changes in minor/patch releases
- When breaking changes are necessary, provide migration path
- Deprecate old APIs before removing them
- Use build tags for experimental features

```go
// Deprecation example
// Deprecated: Use NewConnectionWithConfig instead.
func NewConnection(dsn string) (*Connection, error) {
    return NewConnectionWithConfig(context.Background(), ConfigParams{DSN: dsn})
}
```

### API Evolution Patterns

- Add new optional parameters using functional options
- Extend interfaces by embedding, don't modify existing ones
- Use struct embedding to extend types
- Provide adapter functions for compatibility

## Error Handling in APIs

### Error Types

- Use sentinel errors for expected conditions
- Use custom error types for complex error information
- Wrap errors with context using fmt.Errorf with %w verb
- Provide error inspection functions

```go
// Sentinel errors
var (
    ErrNotFound     = errors.New("item not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrTimeout      = errors.New("operation timed out")
)

// Custom error types
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}

// Error inspection
func IsNotFound(err error) bool {
    return errors.Is(err, ErrNotFound)
}

func IsValidationError(err error) bool {
    var ve ValidationError
    return errors.As(err, &ve)
}
```

### Error Context

- Include relevant context in error messages
- Don't expose internal implementation details
- Use structured errors for programmatic handling
- Log detailed errors internally, return simple errors to users

## Documentation Standards

### Package Documentation

```go
// Package config provides configuration loading with auto-discovery.
//
// The config package supports multiple configuration sources including
// YAML files, environment variables, and command-line flags. It automatically
// discovers configuration files in common locations and merges them according
// to a defined precedence order.
//
// Basic usage:
//
//	type AppConfig struct {
//	    Port int    `yaml:"port"`
//	    Host string `yaml:"host"`
//	}
//
//	cfg := config.MustLoad[AppConfig]()
//
// For more advanced usage, see the examples in the examples/ directory.
package config
```

### Function Documentation

```go
// NewConnection creates a new database connection with the given configuration.
//
// The connection is established immediately and tested with a ping operation.
// If the connection fails, an error is returned with details about the failure.
//
// Options can be provided to customize connection behavior:
//   - WithTimeout: sets connection timeout (default: 30s)
//   - WithRetries: sets number of retry attempts (default: 3)
//   - WithTracing: enables OpenTelemetry tracing (default: true)
//
// Example:
//
//	conn, err := NewConnection(ctx, "postgres://user:pass@localhost/db",
//	    WithTimeout(10*time.Second),
//	    WithRetries(5),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer conn.Close()
func NewConnection(ctx context.Context, dsn string, opts ...Option) (*Connection, error)
```

## Testing Public APIs

### API Contract Testing

- Test all public functions and methods
- Test with various input combinations
- Test error conditions and edge cases
- Test concurrent usage patterns

### Example-Based Testing

```go
func ExampleNewConnection() {
    ctx := context.Background()

    conn, err := NewConnection(ctx, "postgres://localhost/test")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Use connection...
    fmt.Println("Connected successfully")
    // Output: Connected successfully
}
```

### Compatibility Testing

- Test against multiple Go versions
- Test with different dependency versions
- Test upgrade/downgrade scenarios
- Use integration tests for real-world usage

## Performance Considerations

### API Performance

- Design APIs to minimize allocations
- Use object pooling for frequently created objects
- Provide batch operations for bulk operations
- Consider streaming APIs for large datasets

### Resource Management

- Provide explicit cleanup methods (Close, Shutdown)
- Use context for cancellation and timeouts
- Implement proper connection pooling
- Handle resource exhaustion gracefully

```go
// Good: Explicit resource management
type Connection struct {
    // internal fields
}

func (c *Connection) Close() error {
    // cleanup resources
    return nil
}

// Good: Context support
func (c *Connection) Query(ctx context.Context, query string) (*Result, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        // perform query
    }
}
```
