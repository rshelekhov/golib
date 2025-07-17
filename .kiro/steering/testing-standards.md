# Testing Standards and Practices

## Test Organization

### File Structure

- Place tests in `*_test.go` files alongside source code
- Use `testutil/` package for shared testing utilities
- Create `examples/` directory with runnable examples
- Use `integration_test.go` suffix for integration tests

### Test Naming

- Test functions: `TestFunctionName_Scenario`
- Benchmark functions: `BenchmarkFunctionName`
- Example functions: `ExampleFunctionName`
- Helper functions: `testHelperName` (unexported)

## Unit Testing Patterns

### Table-Driven Tests

```go
func TestValidateConfig(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid config",
            config:  Config{ServiceName: "test", Port: 8080},
            wantErr: false,
        },
        {
            name:    "missing service name",
            config:  Config{Port: 8080},
            wantErr: true,
            errMsg:  "service name is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

### Test Assertions

- Use `testify/require` for critical assertions that should stop the test
- Use `testify/assert` for non-critical assertions
- Prefer specific assertion methods over generic ones
- Include meaningful error messages in assertions

### Mocking and Stubs

- Use interfaces to enable mocking
- Create mocks in `testutil/` package when shared
- Use `testify/mock` for complex mocking scenarios
- Keep mocks simple and focused

## Integration Testing

### Docker Container Testing

```go
func TestDatabaseOperations(t *testing.T) {
    ctx := context.Background()

    // Use testutil to create test database
    testDB, err := testutil.NewTestDB(ctx)
    require.NoError(t, err)
    defer testDB.Close(ctx)

    // Create connection using test database
    conn, err := NewConnection(ctx, testDB.ConnectionString())
    require.NoError(t, err)
    defer conn.Close()

    // Run tests against real database
    // ...
}
```

### Environment Variables for Testing

- Support `TEST_*` environment variables for external services
- Fall back to Docker containers when env vars not set
- Document required test environment setup
- Use `testing.Short()` to skip slow tests

### Test Data Management

- Use factories or builders for test data creation
- Clean up test data after each test
- Use transactions that rollback for database tests
- Generate realistic test data with libraries like `gofakeit`
- Follow the testutil patterns defined in golib-architecture.md

## Performance Testing

### Benchmarks

```go
func BenchmarkConnectionPool(b *testing.B) {
    pool := setupConnectionPool()
    defer pool.Close()

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            conn := pool.Get()
            // Perform operation
            pool.Put(conn)
        }
    })
}
```

### Load Testing

- Include load tests for critical paths
- Test connection pool behavior under load
- Verify graceful degradation under stress
- Test timeout and retry mechanisms

## Test Utilities

### Helper Functions

```go
// testutil/database.go
func NewTestDB(ctx context.Context) (*TestDB, error) {
    // Start Docker container or connect to existing DB
    // Return connection details
}

func (db *TestDB) CreateTestData(ctx context.Context) error {
    // Insert known test data
}

func (db *TestDB) Cleanup(ctx context.Context) error {
    // Clean up test data
}
```

### Test Fixtures

- Provide common test configurations
- Include sample data files (JSON, YAML)
- Create reusable test scenarios
- Document test setup requirements

## CI/CD Testing

### Test Categories

- Unit tests: Fast, no external dependencies
- Integration tests: Require external services
- End-to-end tests: Full system testing
- Performance tests: Benchmark and load tests

### Test Execution

- Run unit tests on every commit
- Run integration tests on pull requests
- Run performance tests on releases
- Use build tags to separate test types

### Coverage Requirements

- Aim for 80%+ code coverage on new code
- Require coverage reports in CI
- Focus on testing critical paths thoroughly
- Don't sacrifice test quality for coverage numbers

## Testing Best Practices

### Test Independence

- Each test should be independent and isolated
- Tests should not depend on execution order
- Clean up resources in defer statements
- Use fresh test data for each test

### Error Testing

- Test both success and failure scenarios
- Verify error messages and types
- Test edge cases and boundary conditions
- Test timeout and cancellation scenarios

### Documentation

- Include examples that serve as documentation
- Test examples to ensure they work
- Document test setup requirements
- Explain complex test scenarios in comments
