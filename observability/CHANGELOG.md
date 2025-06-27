# Changelog

All notable changes to the Observability package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
