# Changelog

All notable changes to the Request ID middleware package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-30

### Added

- Initial release of Request ID middleware package
- Extract request ID from gRPC metadata or HTTP headers
- Generate new request ID using ksuid if not present
- Store request ID in context for use throughout request lifecycle
- Support for both unary and streaming gRPC interceptors
- HTTP middleware for standard HTTP handlers
- Context utilities (`FromContext`, `WithContext`)
- Constants for header name and context key
