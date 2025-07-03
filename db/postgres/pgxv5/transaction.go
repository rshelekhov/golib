package pgxv5

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Transaction wraps pgx.Tx to implement QueryEngine interface.
type Transaction struct {
	pgx.Tx
}

// QueryRow executes a query that returns a single row.
func (t *Transaction) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return t.Tx.QueryRow(ctx, sql, args...)
}

// Query executes a query that returns multiple rows.
func (t *Transaction) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return t.Tx.Query(ctx, sql, args...)
}

// Exec executes a query that doesn't return rows.
func (t *Transaction) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return t.Tx.Exec(ctx, sql, args...)
}

// SendBatch sends a batch of queries to the server.
func (t *Transaction) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return t.Tx.SendBatch(ctx, b)
}

// CopyFrom performs a bulk copy operation.
func (t *Transaction) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return t.Tx.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

// Commit commits the transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

// Rollback rolls back the transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}
