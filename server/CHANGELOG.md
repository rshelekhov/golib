# Changelog

All notable changes to the Server Bootstrap package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2025-10-30

### Changed

- Migrated interceptors and middleware to separate middleware packages

### Removed

- Removed `interceptors.go` file (functionality moved to middleware packages)
- Removed `middleware.go` file (functionality moved to middleware packages)

### Added

- Dependency on middleware packages for backward compatibility

## [1.1.0] - 2025-07-11

### Added

- `WithStatsHandler()` option to register a gRPC `stats.Handler` for advanced tracing and metrics collection.
- `ReadinessProvider` interface allowing services to contribute custom readiness checks to the `/readyz` endpoint.

### Changed

- Default README updated to document the new configuration options and readiness extension mechanism.

## [1.0.0] - 2025-06-25

### Added

- Initial release of Server Bootstrap package
- gRPC server with reflection support
- HTTP Gateway server with gRPC-Gateway integration
- Kubernetes-compatible health endpoints (`/healthz`, `/readyz`)
- Graceful shutdown with configurable timeout
- Configurable logging with structured logging support
- Built-in middleware and interceptors:
  - Logging middleware/interceptors
  - Recovery middleware/interceptors
  - CORS middleware
- Functional options pattern for configuration:
  - `WithGRPCPort()` - Configure gRPC server port
  - `WithHTTPPort()` - Configure HTTP server port
  - `WithReflection()` - Enable/disable gRPC reflection
  - `WithShutdownTimeout()` - Configure graceful shutdown timeout
  - `WithUnaryInterceptors()` - Add gRPC unary interceptors
  - `WithStreamInterceptors()` - Add gRPC stream interceptors
  - `WithMuxOptions()` - Add gRPC-Gateway ServeMux options
  - `WithHTTPMiddleware()` - Add HTTP middleware
  - `WithLogger()` - Configure structured logger
- Service interface for easy integration
- Complete example application
- Comprehensive documentation
