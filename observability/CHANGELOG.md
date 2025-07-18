# Changelog

All notable changes to the Observability package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.5.0] - 2025-07-18

### Changed

- **BREAKING**: Changed `ConfigParams.OTLPTransportType` from `tracing.OTLPTransportType` to `string`
  - Provides more flexibility for configuration from environment variables and config files
  - Use `"grpc"` instead of `tracing.OTLPTransportGRPC`
  - Use `"http"` instead of `tracing.OTLPTransportHTTP`
  - Internal `Config.OTLPTransportType` remains typed for type safety
- **Enhanced validation**: Updated validation logic to convert string to typed constant for validation
- **Updated conversion**: Modified `NewConfig` function to handle string-to-type conversion automatically

### Fixed

- **Examples**: Updated all examples to use string values for transport type configuration
- **Tests**: Updated all test cases to use string values instead of typed constants

## [1.4.0] - 2025-07-15

### Added

- **Configurable TLS Support**: Added `OTLPInsecure` configuration option for OTLP exporters
  - Smart environment-based defaults: insecure for local/dev, secure for production
  - Functional option `WithOTLPInsecure(bool)` to override defaults
  - Support for both GRPC and HTTP OTLP transports
  - Applies to all observability components: logging, tracing, and metrics
- **TLS Configuration Examples**: Added comprehensive examples in `examples/tls_example.go`
- **Enhanced Documentation**: Updated README with detailed TLS configuration guide

### Changed

- **OTLP Exporters**: All OTLP exporters now respect the `OTLPInsecure` configuration
  - Logger: `otlploggrpc` exporter uses configurable TLS
  - Tracing: Both `otlptracegrpc` and `otlptracehttp` exporters use configurable TLS
  - Metrics: `otlpmetricgrpc` exporter uses configurable TLS
- **Config Structure**: Added `OTLPInsecure` field to main `Config` and all component configs

### Fixed

- **TLS Consistency**: Resolved hardcoded insecure connections across all OTLP exporters
- **Production Security**: Production environments now use TLS by default for OTLP connections

## [1.3.5] - 2025-07-13

### Added

- **Error helper**: `tracing.RecordError(span, err)` records the error and sets span status to `codes.Error`.

### Documentation

- Updated `observability/tracing/README.md` with error recording examples.

## [1.3.4] - 2025-07-02

### Added

- **Pretty logging for local development**: Colorized, human-readable output for `Env: "local"`
- **PrettyHandler**: New handler with color-coded log levels and formatted timestamps
- Auto-selection of pretty handler for local environment (no OpenTelemetry overhead)

### Changed

- Added `github.com/fatih/color v1.18.0` dependency
- Updated README files and examples

## [1.3.3] - 2025-06-30

### Changed

- **BREAKING**: Improved `NewConfig` API with structured parameters
  - Replaced positional parameters with `ConfigParams` struct for better type safety
  - Added functional options pattern for optional parameters
  - `NewConfig(params ConfigParams, opts ...Option)` instead of `NewConfig(env, serviceName, serviceVersion, enableMetrics, otlpEndpoint, opts...)`
- **Enhanced type safety**: `ConfigParams` struct provides clear field names and compile-time validation

### Documentation

- Updated README.md with new ConfigParams API examples
- Updated examples/main.go to demonstrate new API usage
- Enhanced documentation with clearer parameter descriptions
- Added examples showing both basic and advanced configuration patterns

## [1.3.2] - 2025-06-28

### Changed

- **Dependencies**: Reorganized OpenTelemetry dependencies
  - Moved `go.opentelemetry.io/contrib/bridges/otelslog` to direct dependencies
  - Moved `go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc` to direct dependencies
  - Moved `go.opentelemetry.io/otel/exporters/stdout/stdoutlog` to direct dependencies
  - Moved `go.opentelemetry.io/otel/log` to direct dependencies
  - Moved `go.opentelemetry.io/otel/sdk/log` to direct dependencies
  - Moved `google.golang.org/protobuf` to direct dependencies
- **Removed**: `github.com/fatih/color` dependency (no longer needed)
- **Removed**: `go.opentelemetry.io/otel/exporters/stdout/stdoutmetric` dependency (not used)

## [1.3.1] - 2025-06-28

### Added

- **Unified API**: New `Init(ctx context.Context, cfg Config)` function for all packages
- **Config structs**: Added `Config` type to all packages (logger, metrics, tracing)
- **Automatic exporter selection**: Based on configuration, no manual exporter choice needed
- **Metrics optimization**: Disabled by default in local development environment

### Changed

- Simplified initialization with unified `Init(ctx, Config)` API
- Improved metrics handling in local development (zero overhead)
- Enhanced configuration with automatic exporter selection
- Removed stdout exporter for metrics (impractical)

### Documentation

- Updated all README files with new unified API
- Updated all examples to use new `Init(ctx, Config)` pattern
- Simplified local development setup (no metrics overhead)

## [1.2.0] - 2025-06-27

### Added

- Configurable log levels for all logger initialization functions
- `LogLevel` field to `Config` struct with `slog.Level` type
- **Simplified Configuration API**: `NewConfig(env, serviceName, serviceVersion, enableMetrics, otlpEndpoint, logLevel...)`
- Optional `logLevel` parameter to override environment defaults
- Map-based environment validation for better extensibility
- Methods on Config: `requiresOTLPEndpoint()` and enhanced `Validate()`
- Helper functions: `getDefaultLogLevel(env)` and `getSupportedEnvs()`
- Level filtering in logger handlers for efficient log processing

### Changed

- **BREAKING**: `InitLoggerStdout()` now requires `level slog.Level` parameter
- **BREAKING**: `InitLoggerOTLP()` now requires `level slog.Level` parameter
- **BREAKING**: `NewConfig()` completely redesigned with cleaner, more flexible API
- Renamed `Endpoint` field to `OTLPEndpoint` in `Config` struct for clarity
- Improved error handling with better error messages and centralized validation
- Updated examples to demonstrate new API capabilities
- Improved error handling in `Shutdown()` method using `errors.Join()`

### Documentation

- Updated README.md with simplified configuration examples
- Updated logger/README.md with comprehensive log level documentation
- Added examples showing error handling and configuration flexibility
- Enhanced documentation for different deployment scenarios

## [1.1.0] - 2025-06-25

### Changed

- Updated docker-compose.yml, docker-compose.dev.yml comments
- Updated OpenTelemetry Collector configuration comments (otel-collector-config.yml, otel-collector-dev.yml)
- Updated Prometheus configuration comments (prometheus.yml, prometheus-dev.yml)
- Updated Grafana configuration comments (grafana-dev.yml)
- Improved consistency and international accessibility of documentation

## [1.0.0] - 2025-06-12

### Added

- Initial release of the observability package
- Includes logging functionality (previously logger)
- Added infrastructure for metrics and tracing
- Improved documentation and examples
