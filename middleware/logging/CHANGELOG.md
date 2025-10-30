# Changelog

All notable changes to the Logging middleware package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-30

### Added

- Initial release of Logging middleware package
- gRPC unary interceptor for logging requests with method, status code, and duration
- gRPC stream interceptor for logging streams with method, status code, and duration
- HTTP middleware for logging HTTP requests with method, path, status code, duration, and user agent
- Response writer wrapper to capture HTTP status codes
