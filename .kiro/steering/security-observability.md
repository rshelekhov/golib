# Security and Observability Standards

## Security Best Practices

### Input Validation and Sanitization

- Validate all inputs at API boundaries
- Use allowlists instead of denylists when possible
- Sanitize data before logging to prevent log injection
- Validate configuration parameters thoroughly
- Use type-safe parsing for structured data

```go
// Good: Comprehensive input validation
func (c ConfigParams) Validate() error {
    if c.ServiceName == "" {
        return errors.New("service name is required")
    }
    if !isValidServiceName(c.ServiceName) {
        return errors.New("service name contains invalid characters")
    }
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
    }
    return nil
}

func isValidServiceName(name string) bool {
    // Only allow alphanumeric characters, hyphens, and underscores
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
    return matched
}
```

### Secrets Management

- Never hardcode secrets in source code
- Use environment variables or secret management systems
- Redact secrets in logs and error messages
- Use secure random generation for tokens and IDs
- Implement proper key rotation mechanisms

```go
// Good: Secure configuration handling
type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password" json:"-"` // Exclude from JSON serialization
}

// Redact sensitive information in string representation
func (c DatabaseConfig) String() string {
    return fmt.Sprintf("DatabaseConfig{Host: %s, Port: %d, Username: %s, Password: [REDACTED]}",
        c.Host, c.Port, c.Username)
}

// Use crypto/rand for secure random generation
func generateSecureToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate secure token: %w", err)
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}
```

### TLS and Network Security

- Use TLS for all network communications by default
- Implement proper certificate validation
- Support configurable TLS settings
- Use secure defaults with option to override
- Implement timeout and rate limiting

```go
// Good: Secure TLS configuration
type TLSConfig struct {
    Enabled            bool          `yaml:"enabled"`
    CertFile           string        `yaml:"cert_file"`
    KeyFile            string        `yaml:"key_file"`
    CAFile             string        `yaml:"ca_file"`
    InsecureSkipVerify bool          `yaml:"insecure_skip_verify"`
    MinVersion         uint16        `yaml:"min_version"`
    MaxVersion         uint16        `yaml:"max_version"`
}

func (c TLSConfig) ToGoTLSConfig() (*tls.Config, error) {
    if !c.Enabled {
        return nil, nil
    }

    tlsConfig := &tls.Config{
        MinVersion:         c.MinVersion,
        MaxVersion:         c.MaxVersion,
        InsecureSkipVerify: c.InsecureSkipVerify,
    }

    if c.CertFile != "" && c.KeyFile != "" {
        cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
        if err != nil {
            return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
        }
        tlsConfig.Certificates = []tls.Certificate{cert}
    }

    return tlsConfig, nil
}
```

### Authentication and Authorization

- Implement proper authentication mechanisms
- Use secure session management
- Implement role-based access control where applicable
- Log security events for auditing
- Use constant-time comparison for sensitive data

```go
// Good: Secure token validation
func validateToken(provided, expected string) bool {
    // Use constant-time comparison to prevent timing attacks
    return subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) == 1
}

// Good: Security event logging
func logSecurityEvent(ctx context.Context, event string, details map[string]interface{}) {
    logger := observability.LoggerFromContext(ctx)
    logger.WarnContext(ctx, "security event",
        "event", event,
        "details", details,
        "timestamp", time.Now().UTC(),
    )
}
```

## Observability Implementation

### Security-Focused Logging

- Follow the structured logging standards from golib-architecture.md
- Always sanitize sensitive data before logging to prevent data leaks
- Include security context in log messages (user IDs, IP addresses, etc.)
- Use appropriate log levels for security events (WARN for suspicious activity, ERROR for security failures)

### Security Metrics

- Follow the metrics standards from golib-architecture.md
- Add security-specific metrics: failed authentication attempts, rate limit violations, etc.
- Monitor security-related performance: encryption/decryption times, certificate validation duration
- Track security events for compliance and auditing

### Security Tracing

- Follow the tracing standards from golib-architecture.md
- Include security context in spans (authentication status, authorization levels)
- Sanitize sensitive data in span attributes
- Trace security-critical operations for audit trails
