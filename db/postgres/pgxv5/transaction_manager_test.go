package pgxv5

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rshelekhov/go-db/postgres/pgxv5/testutil"
)

func TestTransactionManager(t *testing.T) {
	ctx := context.Background()

	// Start test database
	db, err := testutil.NewTestDB(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Wait for database to be ready
	err = db.WaitForReady(ctx)
	require.NoError(t, err)

	// Create connection pool
	conn, err := NewConnectionPool(ctx, db.ConnStr())
	require.NoError(t, err)
	defer conn.Close()

	// Create transaction manager
	txManager := NewTransactionManager(conn)

	// Create test table
	_, err = conn.Exec(ctx, `
		CREATE TABLE test (
			id SERIAL PRIMARY KEY,
			value TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	t.Run("Basic Transaction", func(t *testing.T) {
		err := txManager.RunReadCommitted(ctx, func(txCtx context.Context) error {
			_, err := conn.Exec(txCtx, "INSERT INTO test (value) VALUES ($1)", "test1")
			return err
		})
		require.NoError(t, err)

		var value string
		err = conn.QueryRow(ctx, "SELECT value FROM test WHERE value = $1", "test1").Scan(&value)
		require.NoError(t, err)
		assert.Equal(t, "test1", value)
	})

	t.Run("Nested Transaction", func(t *testing.T) {
		err := txManager.RunReadCommitted(ctx, func(txCtx context.Context) error {
			_, err := conn.Exec(txCtx, "INSERT INTO test (value) VALUES ($1)", "test2")
			if err != nil {
				return err
			}

			return txManager.RunReadCommitted(txCtx, func(nestedCtx context.Context) error {
				_, err := conn.Exec(nestedCtx, "INSERT INTO test (value) VALUES ($1)", "test3")
				return err
			})
		})
		require.NoError(t, err)

		var count int
		err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM test WHERE value IN ($1, $2)", "test2", "test3").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("Transaction Rollback", func(t *testing.T) {
		// Clear table before test
		_, err := conn.Exec(ctx, "TRUNCATE TABLE test")
		require.NoError(t, err)

		err = txManager.RunReadCommitted(ctx, func(txCtx context.Context) error {
			_, err := txManager.GetQueryEngine(txCtx).Exec(txCtx, "INSERT INTO test (value) VALUES ($1)", "test4")
			if err != nil {
				return err
			}

			// Force rollback by returning error
			return assert.AnError
		})
		require.Error(t, err)

		var count int
		err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM test WHERE value = $1", "test4").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Read Only Transaction", func(t *testing.T) {
		err := txManager.RunReadCommittedWithAccessMode(ctx, ReadOnly, func(txCtx context.Context) error {
			// Try to insert in read-only transaction
			_, err := txManager.GetQueryEngine(txCtx).Exec(txCtx, "INSERT INTO test (value) VALUES ($1)", "test5")
			if err != nil {
				return err
			}
			return nil
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only transaction")
	})

	t.Run("Serializable Transaction", func(t *testing.T) {
		err := txManager.RunSerializable(ctx, func(txCtx context.Context) error {
			_, err := txManager.GetQueryEngine(txCtx).Exec(txCtx, "INSERT INTO test (value) VALUES ($1)", "test6")
			return err
		})
		require.NoError(t, err)

		var value string
		err = conn.QueryRow(ctx, "SELECT value FROM test WHERE value = $1", "test6").Scan(&value)
		require.NoError(t, err)
		assert.Equal(t, "test6", value)
	})

	t.Run("RepeatableRead Transaction", func(t *testing.T) {
		err := txManager.RunRepeatableRead(ctx, func(txCtx context.Context) error {
			_, err := txManager.GetQueryEngine(txCtx).Exec(txCtx, "INSERT INTO test (value) VALUES ($1)", "test7")
			return err
		})
		require.NoError(t, err)

		var value string
		err = conn.QueryRow(ctx, "SELECT value FROM test WHERE value = $1", "test7").Scan(&value)
		require.NoError(t, err)
		assert.Equal(t, "test7", value)
	})

	t.Run("RepeatableRead With Access Mode", func(t *testing.T) {
		err := txManager.RunRepeatableReadWithAccessMode(ctx, ReadOnly, func(txCtx context.Context) error {
			_, err := txManager.GetQueryEngine(txCtx).Exec(txCtx, "INSERT INTO test (value) VALUES ($1)", "test8")
			if err != nil {
				return err
			}
			return nil
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only transaction")
	})

	t.Run("Serializable With Access Mode", func(t *testing.T) {
		err := txManager.RunSerializableWithAccessMode(ctx, ReadOnly, func(txCtx context.Context) error {
			_, err := txManager.GetQueryEngine(txCtx).Exec(txCtx, "INSERT INTO test (value) VALUES ($1)", "test9")
			if err != nil {
				return err
			}
			return nil
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only transaction")
	})

	t.Run("GetQueryEngine", func(t *testing.T) {
		// Outside transaction should return connection
		engine := txManager.GetQueryEngine(ctx)
		assert.IsType(t, &Connection{}, engine)

		// Inside transaction should return transaction
		err := txManager.RunReadCommitted(ctx, func(txCtx context.Context) error {
			engine := txManager.GetQueryEngine(txCtx)
			assert.IsType(t, &Transaction{}, engine)
			return nil
		})
		require.NoError(t, err)
	})
}
