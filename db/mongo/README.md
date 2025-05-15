# MongoDB wrapper

MongoDB database client with transaction management.

## Features

- Connection management
- Transaction management
- Session handling
- Server API version support
- Interface-based design for better abstraction

## Usage

```go
// Create connection
conn, err := mongo.NewConnection(ctx, "mongodb://localhost:27017", "mydb",
    mongo.WithTimeout(time.Second*5),
    mongo.WithServerAPI("1"),
)
if err != nil {
    log.Fatal(err)
}
defer conn.Close(ctx)

// Create transaction manager
txManager := mongo.NewTransactionManager(conn)

// Run transaction
err = txManager.RunTransaction(ctx, func(ctx context.Context) error {
    // Use ctx for database operations
    db := conn.Database().(*mongo.Database)
    collection := db.Collection("users")
    _, err := collection.InsertOne(ctx, bson.M{"name": "John"})
    return err
})
```

## Transaction Management

MongoDB transactions are handled through sessions. The transaction manager provides a simple interface to execute operations within a transaction:

```go
err = txManager.RunTransaction(ctx, func(ctx context.Context) error {
    // All operations in this function will be executed in a single transaction
    return nil
})
```

## Connection Options

- `WithTimeout` - Sets the connection timeout
- `WithServerAPI` - Sets the server API version
