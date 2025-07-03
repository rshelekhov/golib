package pgxv5

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// TransactionManager manages database transactions with different isolation levels and access modes.
type TransactionManager struct {
	conn *Connection
}

// NewTransactionManager creates a new transaction manager.
func NewTransactionManager(conn *Connection) *TransactionManager {
	return &TransactionManager{conn: conn}
}

// runTransaction executes the given function within a transaction.
// If a transaction already exists in the context, it will be reused.
func (m *TransactionManager) runTransaction(ctx context.Context, txOpts pgx.TxOptions, fn func(ctx context.Context) error) (err error) {
	// If it's nested Transaction, skip initiating a new one and return func(ctx context.Context) error
	if _, ok := ctx.Value(txKey).(*Transaction); ok {
		return fn(ctx)
	}

	var tx *Transaction

	// Begin runTransaction
	pgxTx, err := m.conn.BeginTx(ctx, txOpts)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %v", err)
	}

	tx = &Transaction{Tx: pgxTx}

	// Set txKey to context
	ctx = context.WithValue(ctx, txKey, tx)

	// Set up a defer function for rolling back the runTransaction.
	defer func() {
		// recover from panic
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		// if func(ctx context.Context) error didn't return error - commit
		if err == nil {
			// if commit returns error -> rollback
			err = tx.Commit(ctx)
			if err != nil {
				err = fmt.Errorf("commit failed: %v", err)
			}
		}

		// rollback on any error
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = fmt.Errorf("rollback failed: %v", errRollback)
			}
		}
	}()

	// Execute the code inside the runTransaction.
	// If the function fails, return the error and the defer function
	//will roll back or commit otherwise.

	// return error without wrapping errors.Wrap
	err = fn(ctx)

	return err
}

// GetQueryEngine returns the appropriate query engine based on the context.
// If a transaction exists in the context, it returns the transaction.
// Otherwise, it returns the connection.
func (m *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	// Transaction always runs on node with NodeRoleWrite role
	if tx, ok := ctx.Value(txKey).(*Transaction); ok {
		return tx
	}

	return m.conn
}

// RunReadCommitted executes the given function within a ReadCommitted transaction.
func (m *TransactionManager) RunReadCommitted(ctx context.Context, f func(txCtx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	}, f)
}

// RunRepeatableRead executes the given function within a RepeatableRead transaction.
func (m *TransactionManager) RunRepeatableRead(ctx context.Context, f func(txCtx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	}, f)
}

// RunSerializable executes the given function within a Serializable transaction.
func (m *TransactionManager) RunSerializable(ctx context.Context, f func(txCtx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	}, f)
}

// RunReadCommittedWithAccessMode executes the given function within a ReadCommitted transaction with specified access mode.
func (m *TransactionManager) RunReadCommittedWithAccessMode(ctx context.Context, accessMode TxAccessMode, f func(txCtx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: accessMode,
	}, f)
}

// RunRepeatableReadWithAccessMode executes the given function within a RepeatableRead transaction with specified access mode.
func (m *TransactionManager) RunRepeatableReadWithAccessMode(ctx context.Context, accessMode TxAccessMode, f func(txCtx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: accessMode,
	}, f)
}

// RunSerializableWithAccessMode executes the given function within a Serializable transaction with specified access mode.
func (m *TransactionManager) RunSerializableWithAccessMode(ctx context.Context, accessMode TxAccessMode, f func(txCtx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: accessMode,
	}, f)
}
