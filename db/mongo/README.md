# MongoDB wrapper

MongoDB database client with transaction management and OpenTelemetry tracing support.

## Features

- Connection management
- Transaction management
- Session handling
- Server API version support
- OpenTelemetry tracing integration (otelmongo)
- Interface-based design for better abstraction

## Usage

```go
// Create connection with tracing enabled (default)
conn, err := mongo.NewConnection(ctx, "mongodb://localhost:27017", "mydb",
    mongo.WithTimeout(time.Second*5),
    mongo.WithServerAPI("1"),
    mongo.WithTracing(true), // can be omitted, default is true
)
if err != nil {
    log.Fatal(err)
}
defer conn.Close(ctx)

// Create connection with tracing disabled
conn, err := mongo.NewConnection(ctx, "mongodb://localhost:27017", "mydb",
    mongo.WithTracing(false), // disable tracing
)

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

- `WithTimeout(duration)` - Sets the connection timeout
- `WithServerAPI(version)` - Sets the server API version
- `WithTracing(bool)` - Enables/disables OpenTelemetry tracing (default: true)
