# Changelog

All notable changes to the MongoDB package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.2] - 2025-07-04

### Changed

- **BREAKING**: Improved type safety for Database and Client methods
  - Changed `Database()` return type from `any` to `*mongo.Database`
  - Changed `Client()` return type from `any` to `*mongo.Client`
  - Updated `ConnectionCloser` interface to reflect specific mongo types

## [1.0.1] - 2025-06-12

### Added

- **OpenTelemetry tracing support**: Added otelmongo integration
- **New configuration option**: `WithTracing(bool)` to control tracing
- **Improved dependency management**: Updated mongo-driver, opentelemetry, snappy, compress dependencies

### Documentation

- Enhanced documentation and usage examples
- Added tracing configuration examples

## [1.0.0] - 2025-05-15

### Added

- **Initial release** of the mongo package
- **Transaction management** and session handling
- **Server API version support** with `WithServerAPI()` option
- **Interface-based design** for better abstraction and testing
- **Connection management** with timeout and URI configuration
- **CRUD operations** with proper error handling
- **Aggregation support** for complex queries
- **Document counting** functionality
