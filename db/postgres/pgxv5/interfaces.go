package pgxv5

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	CommonAPI interface {
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
		Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	}

	TransactionAPI interface {
		BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
		Begin(ctx context.Context) (pgx.Tx, error)
	}

	ExtendedAPI interface {
		SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
		CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	}

	ConnectionAPI interface {
		CommonAPI
		TransactionAPI
		ExtendedAPI
	}

	QueryEngine interface {
		CommonAPI
		ExtendedAPI
	}

	TransactionManagerAPI interface {
		GetQueryEngine(ctx context.Context) QueryEngine
	}
)
