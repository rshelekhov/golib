# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-07-09

### Added

- `WithMinIOEndpoint()` convenience function for MinIO endpoints
- Automatic path-style addressing detection for MinIO endpoints
- `IsMinIOEndpoint()` helper function for endpoint detection
- Auto-SSL disabling for local MinIO endpoints without scheme

### Changed

- Enhanced MinIO compatibility with automatic configuration
- Improved endpoint handling for S3-compatible services

## [1.0.0] - 2025-07-05

### Added

- Initial S3 library implementation
- Connection management with configurable options
- Object operations (Put, Get, Delete, List, Head, Copy)
- Bucket operations (Create, Delete, List, Head, GetLocation)
- Multipart upload support
- Presigned URL generation
- Helper functions for simple operations
- OpenTelemetry tracing integration
- Testing utilities with MinIO container support
- Comprehensive error handling
- AWS credential chain support
- Path-style and virtual-hosted-style addressing
- SSL/TLS configuration options
- Timeout and retry configuration
