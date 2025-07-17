# Go Code Standards and Quality

## Code Style Guidelines

### Formatting and Structure

- Use `gofmt` and `goimports` for consistent formatting
- Line length should not exceed 120 characters
- Use meaningful variable and function names
- Group imports: standard library, third-party, local packages
- Add blank lines between logical sections

### Function Design

- Keep functions small and focused (max 50 lines when possible)
- Use early returns to reduce nesting
- Prefer composition over inheritance
- Use interfaces for abstraction, not concrete types

### Variable Naming

- Use camelCase for unexported identifiers
- Use PascalCase for exported identifiers
- Use short names for short-lived variables (i, j for loops)
- Use descriptive names for longer-lived variables
- Avoid abbreviations unless they're well-known (ctx, cfg, db)

### Comments and Documentation

- Write package-level documentation for all packages
- Document all exported functions, types, and constants
- Use complete sentences in comments
- Start comments with the name of the item being documented
- Add TODO comments for known issues or future improvements

### Error Handling Patterns

```go
// Prefer this pattern
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Use sentinel errors for expected conditions
var ErrNotFound = errors.New("item not found")

// Use custom error types for complex error handling
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}
```

## Code Quality Standards

### Performance Considerations

- Avoid premature optimization, but be aware of common pitfalls
- Use string builders for string concatenation in loops
- Prefer slices over arrays when size is not fixed
- Use sync.Pool for frequently allocated objects
- Profile code when performance is critical

### Memory Management

- Avoid memory leaks by properly closing resources
- Use context.WithTimeout for operations that might hang
- Be careful with goroutine lifecycles
- Use buffered channels appropriately

### Concurrency Patterns

- Use channels for communication, mutexes for state protection
- Prefer select statements over blocking channel operations
- Always handle goroutine cleanup and cancellation
- Use sync.WaitGroup for coordinating goroutines

### Security Best Practices

- Validate all inputs at boundaries
- Use crypto/rand for random number generation
- Sanitize data before logging
- Use TLS for network communications
- Never hardcode secrets in source code
