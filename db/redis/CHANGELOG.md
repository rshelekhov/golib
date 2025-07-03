# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-07-03

### Added

- Initial release of Redis library wrapper
- Connection management with configurable options
- Support for all major Redis data types:
  - String operations (Set, Get, Del, Exists, Expire, TTL)
  - Hash operations (HSet, HGet, HGetAll, HDel, HExists, HKeys, HVals, HLen)
  - List operations (LPush, RPush, LPop, RPop, LLen, LRange)
  - Set operations (SAdd, SRem, SMembers, SIsMember, SCard)
  - Sorted Set operations (ZAdd, ZRem, ZScore, ZRange, ZRevRange, ZCard)
  - Scan operations (Scan, HScan, SScan, ZScan)
- Transaction support using Redis MULTI/EXEC through pipelines
- Pipeline support for batch operations
- OpenTelemetry integration for distributed tracing
- Comprehensive interface-based design for easy testing and mocking
- Docker-based test utilities using testcontainers
- Connection options with sensible defaults:
  - Host/Port configuration
  - Password authentication
  - Database selection
  - Connection pooling (pool size, min idle connections)
  - Timeout configuration (dial, read, write, idle)
  - Retry configuration
  - Tracing enable/disable
- Transaction manager with context-based query engine selection
- Consistent error handling and formatting
- Full test coverage for all operations
