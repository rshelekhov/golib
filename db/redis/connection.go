package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Connection represents a connection to Redis.
type Connection struct {
	client *redis.Client
	tracer trace.Tracer
}

// connectionOptions holds configuration for Redis connection
type connectionOptions struct {
	host          string
	port          int
	password      string
	db            int
	poolSize      int
	minIdleConns  int
	maxRetries    int
	dialTimeout   time.Duration
	readTimeout   time.Duration
	writeTimeout  time.Duration
	idleTimeout   time.Duration
	enableTracing bool
}

// ConnectionOption is a function that configures connection options.
type ConnectionOption func(opts *connectionOptions)

// WithHost sets the Redis host.
func WithHost(host string) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.host = host
	}
}

// WithPort sets the Redis port.
func WithPort(port int) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.port = port
	}
}

// WithPassword sets the Redis password.
func WithPassword(password string) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.password = password
	}
}

// WithDB sets the Redis database number.
func WithDB(db int) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.db = db
	}
}

// WithPoolSize sets the connection pool size.
func WithPoolSize(size int) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.poolSize = size
	}
}

// WithMinIdleConns sets the minimum number of idle connections.
func WithMinIdleConns(conns int) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.minIdleConns = conns
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(retries int) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.maxRetries = retries
	}
}

// WithDialTimeout sets the dial timeout.
func WithDialTimeout(timeout time.Duration) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.dialTimeout = timeout
	}
}

// WithReadTimeout sets the read timeout.
func WithReadTimeout(timeout time.Duration) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.readTimeout = timeout
	}
}

// WithWriteTimeout sets the write timeout.
func WithWriteTimeout(timeout time.Duration) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.writeTimeout = timeout
	}
}

// WithIdleTimeout sets the idle timeout.
func WithIdleTimeout(timeout time.Duration) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.idleTimeout = timeout
	}
}

// WithTracing turns on/off tracing through OpenTelemetry
func WithTracing(enable bool) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.enableTracing = enable
	}
}

// NewConnection creates a new connection to Redis.
func NewConnection(ctx context.Context, opts ...ConnectionOption) (ConnectionAPI, error) {
	// Apply default options
	connOpts := &connectionOptions{
		host:          "localhost",
		port:          6379,
		db:            DefaultDB,
		poolSize:      DefaultPoolSize,
		minIdleConns:  DefaultMinIdleConns,
		maxRetries:    DefaultMaxRetries,
		dialTimeout:   DefaultConnectionTimeout,
		readTimeout:   DefaultConnectionTimeout,
		writeTimeout:  DefaultConnectionTimeout,
		idleTimeout:   DefaultIdleTimeout,
		enableTracing: true, // default is true
	}

	for _, opt := range opts {
		if opt != nil {
			opt(connOpts)
		}
	}

	// Create Redis client options
	clientOpts := &redis.Options{
		Addr:            fmt.Sprintf("%s:%d", connOpts.host, connOpts.port),
		Password:        connOpts.password,
		DB:              connOpts.db,
		PoolSize:        connOpts.poolSize,
		MinIdleConns:    connOpts.minIdleConns,
		MaxRetries:      connOpts.maxRetries,
		DialTimeout:     connOpts.dialTimeout,
		ReadTimeout:     connOpts.readTimeout,
		WriteTimeout:    connOpts.writeTimeout,
		ConnMaxIdleTime: connOpts.idleTimeout,
	}

	client := redis.NewClient(clientOpts)

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	conn := &Connection{
		client: client,
	}

	if connOpts.enableTracing {
		conn.tracer = otel.Tracer("redis")
	}

	return conn, nil
}

// Close closes the connection to Redis.
func (c *Connection) Close() error {
	return c.client.Close()
}

// Client returns the Redis client.
func (c *Connection) Client() *redis.Client {
	return c.client
}

// Ping checks the connection to the Redis server.
func (c *Connection) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// String operations
func (c *Connection) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Connection) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Connection) Del(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Del(ctx, keys...).Result()
}

func (c *Connection) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

func (c *Connection) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

func (c *Connection) ExpireAt(ctx context.Context, key string, tm time.Time) error {
	return c.client.ExpireAt(ctx, key, tm).Err()
}

func (c *Connection) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// Hash operations
func (c *Connection) HSet(ctx context.Context, key string, values ...any) error {
	return c.client.HSet(ctx, key, values...).Err()
}

func (c *Connection) HGet(ctx context.Context, key, field string) (string, error) {
	return c.client.HGet(ctx, key, field).Result()
}

func (c *Connection) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

func (c *Connection) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return c.client.HDel(ctx, key, fields...).Result()
}

func (c *Connection) HExists(ctx context.Context, key, field string) (bool, error) {
	return c.client.HExists(ctx, key, field).Result()
}

func (c *Connection) HKeys(ctx context.Context, key string) ([]string, error) {
	return c.client.HKeys(ctx, key).Result()
}

func (c *Connection) HVals(ctx context.Context, key string) ([]string, error) {
	return c.client.HVals(ctx, key).Result()
}

func (c *Connection) HLen(ctx context.Context, key string) (int64, error) {
	return c.client.HLen(ctx, key).Result()
}

// List operations
func (c *Connection) LPush(ctx context.Context, key string, values ...any) (int64, error) {
	return c.client.LPush(ctx, key, values...).Result()
}

func (c *Connection) RPush(ctx context.Context, key string, values ...any) (int64, error) {
	return c.client.RPush(ctx, key, values...).Result()
}

func (c *Connection) LPop(ctx context.Context, key string) (string, error) {
	return c.client.LPop(ctx, key).Result()
}

func (c *Connection) RPop(ctx context.Context, key string) (string, error) {
	return c.client.RPop(ctx, key).Result()
}

func (c *Connection) LLen(ctx context.Context, key string) (int64, error) {
	return c.client.LLen(ctx, key).Result()
}

func (c *Connection) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.LRange(ctx, key, start, stop).Result()
}

// Set operations
func (c *Connection) SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return c.client.SAdd(ctx, key, members...).Result()
}

func (c *Connection) SRem(ctx context.Context, key string, members ...any) (int64, error) {
	return c.client.SRem(ctx, key, members...).Result()
}

func (c *Connection) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

func (c *Connection) SIsMember(ctx context.Context, key string, member any) (bool, error) {
	return c.client.SIsMember(ctx, key, member).Result()
}

func (c *Connection) SCard(ctx context.Context, key string) (int64, error) {
	return c.client.SCard(ctx, key).Result()
}

// Sorted Set operations
func (c *Connection) ZAdd(ctx context.Context, key string, members ...redis.Z) (int64, error) {
	return c.client.ZAdd(ctx, key, members...).Result()
}

func (c *Connection) ZRem(ctx context.Context, key string, members ...any) (int64, error) {
	return c.client.ZRem(ctx, key, members...).Result()
}

func (c *Connection) ZScore(ctx context.Context, key, member string) (float64, error) {
	return c.client.ZScore(ctx, key, member).Result()
}

func (c *Connection) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

func (c *Connection) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRevRange(ctx, key, start, stop).Result()
}

func (c *Connection) ZCard(ctx context.Context, key string) (int64, error) {
	return c.client.ZCard(ctx, key).Result()
}

// Scan operations
func (c *Connection) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.client.Scan(ctx, cursor, match, count).Result()
}

func (c *Connection) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.client.HScan(ctx, key, cursor, match, count).Result()
}

func (c *Connection) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.client.SScan(ctx, key, cursor, match, count).Result()
}

func (c *Connection) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.client.ZScan(ctx, key, cursor, match, count).Result()
}

// Pipeline operations
func (c *Connection) Pipeline() redis.Pipeliner {
	return c.client.Pipeline()
}

func (c *Connection) TxPipeline() redis.Pipeliner {
	return c.client.TxPipeline()
}
