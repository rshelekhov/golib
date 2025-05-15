package pgxv5

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type connectionPoolOptions struct {
	maxConnIdleTime     time.Duration
	maxConnLifeTime     time.Duration
	minConnectionsCount int32
	maxConnectionsCount int32
	tlsConfig           *tls.Config
}

// ConnectionPoolOption is a function that configures connection pool options.
type ConnectionPoolOption func(options *connectionPoolOptions)

// WithMaxConnIdleTime sets the maximum amount of time a connection can be idle.
func WithMaxConnIdleTime(d time.Duration) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.maxConnIdleTime = d
	}
}

// WithMaxConnLifeTime sets the maximum amount of time a connection can exist.
func WithMaxConnLifeTime(d time.Duration) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.maxConnLifeTime = d
	}
}

// WithMinConnectionsCount sets the minimum number of connections in the pool.
func WithMinConnectionsCount(c int32) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.minConnectionsCount = c
	}
}

// WithMaxConnectionsCount sets the maximum number of connections in the pool.
func WithMaxConnectionsCount(c int32) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.maxConnectionsCount = c
	}
}

// WithTLS sets the TLS configuration for the connection.
func WithTLS(cfg *tls.Config) ConnectionPoolOption {
	return func(opts *connectionPoolOptions) {
		opts.tlsConfig = cfg
	}
}

// Connection represents a connection pool to the database.
type Connection struct {
	pool *pgxpool.Pool
}

var (
	_ CommonAPI      = (*Connection)(nil)
	_ ExtendedAPI    = (*Connection)(nil)
	_ TransactionAPI = (*Connection)(nil)
)

// NewConnectionPool creates a new connection pool with the given connection string and options.
func NewConnectionPool(ctx context.Context, connString string, opts ...ConnectionPoolOption) (*Connection, error) {
	// parse connString
	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("can't parse connection string to config: %w", err)
	}

	// ...
	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	// make options
	options := &connectionPoolOptions{
		maxConnIdleTime:     maxConnIdleTimeDefault,
		maxConnLifeTime:     maxConnLifeTimeDefault,
		minConnectionsCount: minConnectionsCountDefault,
		maxConnectionsCount: maxConnectionsCountDefault,
	}
	for _, opt := range opts {
		opt(options)
	}

	// apply options
	connConfig.MaxConnIdleTime = options.maxConnIdleTime
	connConfig.MaxConnLifetime = options.maxConnLifeTime
	connConfig.MinConns = options.minConnectionsCount
	connConfig.MaxConns = options.maxConnectionsCount
	connConfig.ConnConfig.Config.TLSConfig = options.tlsConfig

	// connect to database
	p, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	// ping database
	if err := p.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database error: %w", err)
	}

	return &Connection{
		pool: p,
	}, nil
}

// Close closes the connection pool.
func (c *Connection) Close() {
	c.pool.Close()
}

// Query executes a query that returns multiple rows.
func (c *Connection) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that returns a single row.
func (c *Connection) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}

// Exec executes a query that doesn't return rows.
func (c *Connection) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, args...)
}

// Begin starts a new transaction.
func (c *Connection) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.pool.Begin(ctx)
}

// BeginTx starts a new transaction with the given options.
func (c *Connection) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, txOptions)
}

// SendBatch sends a batch of queries to the server.
func (c *Connection) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return c.pool.SendBatch(ctx, b)
}

// CopyFrom performs a bulk copy operation.
func (c *Connection) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return c.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

// Pool returns the underlying connection pool.
func (c *Connection) Pool() *pgxpool.Pool {
	return c.pool
}
