# Postgres pgx wrapper

PostgreSQL database client with transaction management.

## Features

- Connection pool management
- Transaction management with different isolation levels
- Support for read-only transactions
- Batch operations
- COPY operations
- UUID support

## Usage

```go
// Create connection pool
conn, err := pgx.NewConnectionPool(ctx, "postgres://user:pass@localhost:5432/dbname")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Create transaction manager
txManager := pgx.NewTransactionManager(conn)

// Run transaction with ReadCommitted isolation level
err = txManager.RunReadCommitted(ctx, func(txCtx context.Context) error {
    // Use txManager.GetQueryEngine(txCtx) for database operations
    _, err := txManager.GetQueryEngine(txCtx).Exec(txCtx, "INSERT INTO users (name) VALUES ($1)", "John")
    return err
})

// Run read-only transaction
err = txManager.RunReadCommittedWithAccessMode(ctx, pgx.ReadOnly, func(txCtx context.Context) error {
    // Use txManager.GetQueryEngine(txCtx) for read-only operations
    var name string
    err := txManager.GetQueryEngine(txCtx).QueryRow(txCtx, "SELECT name FROM users WHERE id = $1", 1).Scan(&name)
    return err
})

// Run nested transaction
err = txManager.RunReadCommitted(ctx, func(txCtx context.Context) error {
    // First operation in transaction
    _, err := txManager.GetQueryEngine(txCtx).Exec(txCtx, "INSERT INTO users (name) VALUES ($1)", "John")
    if err != nil {
        return err
    }

    // Nested transaction reuses parent transaction
    return txManager.RunReadCommitted(txCtx, func(nestedCtx context.Context) error {
        _, err := txManager.GetQueryEngine(nestedCtx).Exec(nestedCtx, "INSERT INTO profiles (user_id) VALUES ($1)", 1)
        return err
    })
})
```

## Transaction Isolation Levels

- `ReadCommitted` - Default PostgreSQL isolation level
- `RepeatableRead` - Serializable snapshot isolation
- `Serializable` - True serializable isolation

## Access Modes

- `ReadWrite` - Default mode, allows both reads and writes
- `ReadOnly` - Read-only mode, prevents any modifications
