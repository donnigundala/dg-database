# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-27

### Added
- Initial stable release of the `dg-database` plugin.
- **Multi-Driver Support**: MySQL, PostgreSQL, and SQLite with unified API.
- **Connection Management**: Pooling, read/write splitting, and multi-connection support.
- **Load Balancing**: Round-robin, random, and weighted strategies for read replicas.
- **Transaction Support**: Full ACID transactions with context and savepoint support.
- **Migration System**: Version-controlled schema changes with up/down/reset functionality.
- **PostgreSQL Schema Support**: Multi-tenancy via `search_path` parameter.
- **Container Integration**: Auto-registration of named connections with Injectable pattern.
- **Helper Functions**: `Resolve()`, `MustResolve()`, `ResolveConnection()`, `MustResolveConnection()`.
- **Health Monitoring**: Connection health checks across all databases.
- **Observability**: OpenTelemetry metrics for connection pool and query performance.

### Features
- Thread-safe connection management
- Automatic routing of reads to slaves, writes to master
- Runtime connection add/remove capabilities
- Fluent configuration API with builder pattern
- Auto-migration support for rapid development
- Service provider integration with dg-core framework
- Comprehensive test coverage (40+ tests)

### Documentation
- Complete API reference (documentation/API.md)
- Migration guides (documentation/MIGRATIONS.md, documentation/GO_MIGRATE.md)
- PostgreSQL schema support guide (documentation/SCHEMA.md)
- 6 working examples with READMEs

### Performance
- Connection pooling with configurable limits
- Read replica load balancing
- Efficient transaction management

---

## Development History

The following versions represent the development journey leading to v1.0.0:

### 2025-12-05
- Added DB interface for better abstraction
- Auto-registration of named connections in container
- Injectable helper for easier dependency injection

### 2025-11-24
- PostgreSQL schema support via `search_path`
- Schema-based multi-tenancy patterns
- Comprehensive schema documentation

### 2025-11-23
- Initial beta release with core functionality
- MySQL, PostgreSQL, and SQLite support
- Read/write splitting and migrations
