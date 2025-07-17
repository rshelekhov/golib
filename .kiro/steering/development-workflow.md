# Development Workflow and CI/CD Standards

## Development Environment Setup

### Local Development Requirements

- Go 1.24.2 or later (as specified in go.work)
- Docker and Docker Compose for integration testing
- golangci-lint for code quality checks
- gofmt and goimports for code formatting
- Make for build automation

### IDE Configuration

- Configure gofmt and goimports to run on save
- Set up golangci-lint integration
- Configure test runner for table-driven tests
- Set up debugging for integration tests with Docker

### Environment Variables

```bash
# Development environment
export GO_ENV=development
export LOG_LEVEL=debug

# Testing environment
export TEST_TIMEOUT=30s
export TEST_INTEGRATION=true
export TEST_DOCKER_CLEANUP=true
```

## Git Workflow Standards

### Branch Naming

- `feature/description` - New features
- `fix/description` - Bug fixes
- `refactor/description` - Code refactoring
- `docs/description` - Documentation updates
- `test/description` - Test improvements

### Commit Message Format

```
type(scope): brief description

Detailed explanation of the change, including:
- What was changed and why
- Any breaking changes
- References to issues or tickets

Closes #123
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `perf`, `style`

### Pull Request Guidelines

- Include comprehensive description of changes
- Reference related issues or tickets
- Include test coverage information
- Update documentation if needed
- Ensure all CI checks pass

## Code Quality Standards

### Pre-commit Checks

```bash
# Format code
gofmt -w .
goimports -w .

# Run linter
golangci-lint run

# Run tests
go test ./...

# Check for vulnerabilities
govulncheck ./...

# Check dependencies
go mod tidy
go mod verify
```

### Automated Quality Gates

- Code coverage must be >= 80% for new code
- All linter checks must pass
- No security vulnerabilities allowed
- All tests must pass
- Documentation must be updated

### Code Review Checklist

- [ ] Code follows Go idioms and conventions
- [ ] Error handling is comprehensive and appropriate
- [ ] Tests cover happy path and error cases
- [ ] Documentation is clear and complete
- [ ] No hardcoded secrets or sensitive data
- [ ] Performance implications considered
- [ ] Security implications reviewed
- [ ] Backward compatibility maintained

## Build and Release Process

### Makefile Standards

```makefile
.PHONY: help build test lint clean install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all modules
	go build ./...

test: ## Run all tests
	go test -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests
	go test -tags=integration -race ./...

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	gofmt -w .
	goimports -w .

clean: ## Clean build artifacts
	go clean ./...
	rm -f coverage.out

install: ## Install dependencies
	go mod download
	go mod tidy

security: ## Run security checks
	govulncheck ./...

coverage: test ## Generate coverage report
	go tool cover -html=coverage.out -o coverage.html
```

### Version Management

- Use semantic versioning (semver) for all modules
- Tag releases with `v` prefix (e.g., `v1.2.3`)
- Maintain CHANGELOG.md for each module
- Use Go modules for dependency management

### Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version bumped appropriately
- [ ] Security scan completed
- [ ] Performance benchmarks run
- [ ] Breaking changes documented
- [ ] Migration guide provided (if needed)

### Deployment Standards

- Use Docker containers for consistent deployments
- Implement health checks in all services
- Use configuration management for environment-specific settings
- Implement graceful shutdown handling
- Monitor deployment success and rollback if needed

## Documentation Standards

### Documentation Standards

- Follow the README structure defined in golib-architecture.md
- Use the documentation patterns from api-design-standards.md
- Generate API documentation with `godoc`
- Keep documentation in sync with code changes

## Development Observability

- Follow observability standards from golib-architecture.md and security-observability.md
- Use debug logging during development for troubleshooting
- Profile performance-critical code paths during development
- Implement debug endpoints for local troubleshooting

## Dependency Management

### Dependency Selection Criteria

- Actively maintained projects
- Good security track record
- Stable API with semantic versioning
- Comprehensive documentation
- Compatible license

### Dependency Updates

- Regular security updates
- Automated dependency scanning
- Test compatibility before updating
- Document breaking changes
- Maintain dependency inventory

### Vendor Management

- Use `go mod vendor` for critical dependencies
- Regular vendor directory updates
- Security scanning of vendored code
- License compliance checking

```

```
