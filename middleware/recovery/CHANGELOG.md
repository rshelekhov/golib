# Changelog

All notable changes to the Recovery middleware package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-30

### Added

- Initial release of Recovery middleware package
- gRPC unary interceptor for recovering from panics
- gRPC stream interceptor for recovering from panics
- HTTP middleware for recovering from panics
- Panic logging with error details and request context
- Automatic error response generation (Internal Server Error for gRPC, HTTP 500 for HTTP)
