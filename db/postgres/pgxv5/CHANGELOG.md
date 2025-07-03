# Changelog

All notable changes to the pgxv5 PostgreSQL adapter will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-07-03

### Changed

- **BREAKING**: Moved package to `pgxv5/` subdirectory for better organization
- **BREAKING**: Module path changed from `github.com/rshelekhov/golib/db/postgres` to `github.com/rshelekhov/golib/db/postgres/pgxv5`
- Enhanced package structure as part of modular PostgreSQL adapters architecture

All functionality remains identical.

## [1.0.1] - 2025-06-12

### Added

- OpenTelemetry tracing support (otelpgx)
- New `WithTracing(bool)` option to control tracing

### Changed

- Connection pool creation (new parameters)
- Updated dependencies (added otelpgx)
- Improved documentation and usage examples

## [1.0.0] - 2025-01-05

### Added

- Initial release of the postgres package
- Transaction management and session handling
- Connection pool with custom options
- Interface-based design for abstraction
