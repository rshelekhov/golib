package redis

import "time"

const (
	// DefaultConnectionTimeout is the default timeout for Redis connection
	DefaultConnectionTimeout = 10 * time.Second
	// DefaultIdleTimeout is the default timeout for idle connections
	DefaultIdleTimeout = 5 * time.Minute
	// DefaultMaxRetries is the default number of retries for Redis operations
	DefaultMaxRetries = 3
	// DefaultPoolSize is the default size of the connection pool
	DefaultPoolSize = 10
	// DefaultMinIdleConns is the default minimum number of idle connections
	DefaultMinIdleConns = 5
	// DefaultDB is the default database number
	DefaultDB = 0
)
