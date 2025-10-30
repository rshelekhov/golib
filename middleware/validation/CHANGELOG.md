# Changelog

All notable changes to the Validation middleware package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-30

### Added

- Initial release of Validation middleware package
- gRPC unary interceptor for validating requests
- Automatic validation of requests implementing `Validate() error` method
- Returns InvalidArgument status on validation failure
