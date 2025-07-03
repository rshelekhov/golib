# Redis Library

A Redis wrapper library that provides a unified interface for Redis operations with support for pipelines, transactions, and connection management.

## Features

- **Connection Management**: Easy Redis connection setup with configurable options
- **Transaction Support**: Redis transactions using MULTI/EXEC through pipelines
- **Pipeline Support**: Batch operations for improved performance
- **Comprehensive API**: Support for all major Redis data types (strings, hashes, lists, sets, sorted sets)
- **Scan Operations**: Efficient iteration over large datasets
- **OpenTelemetry Integration**: Built-in tracing support
- **Testing Utilities**: Docker-based test utilities for integration testing

## Installation

```bash
go get github.com/rshelekhov/golib/db/redis
```

## Usage

### Basic Connection

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/rshelekhov/golib/db/redis"
)

func main() {
    ctx := context.Background()

    // Create connection with default options
    conn, err := redis.NewConnection(ctx)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    // Set a key
    err = conn.Set(ctx, "key", "value", time.Hour)
    if err != nil {
        panic(err)
    }

    // Get a key
    value, err := conn.Get(ctx, "key")
    if err != nil {
        panic(err)
    }
    fmt.Println(value) // Output: value
}
```

### Connection with Custom Options

```go
conn, err := redis.NewConnection(ctx,
    redis.WithHost("localhost"),
    redis.WithPort(6379),
    redis.WithPassword("password"),
    redis.WithDB(0),
    redis.WithPoolSize(20),
    redis.WithMinIdleConns(10),
    redis.WithMaxRetries(3),
    redis.WithDialTimeout(5*time.Second),
    redis.WithReadTimeout(3*time.Second),
    redis.WithWriteTimeout(3*time.Second),
    redis.WithIdleTimeout(5*time.Minute),
    redis.WithTracing(true),
)
```

### Transaction Support

```go
// Create transaction manager
tm := redis.NewTransactionManager(conn.(*redis.Connection))

// Run transaction
err := tm.RunTransaction(ctx, func(ctx context.Context) error {
    engine := tm.GetQueryEngine(ctx)

    // All operations will be queued in the transaction
    if err := engine.Set(ctx, "key1", "value1", 0); err != nil {
        return err
    }

    if err := engine.Set(ctx, "key2", "value2", 0); err != nil {
        return err
    }

    // Transaction will be committed automatically
    return nil
})
```

### Pipeline Support

```go
// Run pipeline (non-transactional batching)
err := tm.RunPipeline(ctx, func(ctx context.Context) error {
    engine := tm.GetQueryEngine(ctx)

    // All operations will be batched
    if err := engine.Set(ctx, "key1", "value1", 0); err != nil {
        return err
    }

    if err := engine.Set(ctx, "key2", "value2", 0); err != nil {
        return err
    }

    // Pipeline will be executed automatically
    return nil
})
```

### Hash Operations

```go
// Set hash fields
err := conn.HSet(ctx, "user:123", "name", "John", "age", "30")
if err != nil {
    panic(err)
}

// Get hash field
name, err := conn.HGet(ctx, "user:123", "name")
if err != nil {
    panic(err)
}

// Get all hash fields
user, err := conn.HGetAll(ctx, "user:123")
if err != nil {
    panic(err)
}
fmt.Println(user) // map[name:John age:30]
```

### List Operations

```go
// Push to list
count, err := conn.LPush(ctx, "mylist", "item1", "item2", "item3")
if err != nil {
    panic(err)
}

// Get list range
items, err := conn.LRange(ctx, "mylist", 0, -1)
if err != nil {
    panic(err)
}
fmt.Println(items) // [item3 item2 item1]
```

### Set Operations

```go
// Add to set
count, err := conn.SAdd(ctx, "myset", "member1", "member2", "member3")
if err != nil {
    panic(err)
}

// Get all members
members, err := conn.SMembers(ctx, "myset")
if err != nil {
    panic(err)
}
fmt.Println(members) // [member1 member2 member3]
```

### Sorted Set Operations

```go
// Add to sorted set
err := conn.ZAdd(ctx, "leaderboard",
    redis.Z{Score: 100, Member: "player1"},
    redis.Z{Score: 200, Member: "player2"},
    redis.Z{Score: 150, Member: "player3"},
)
if err != nil {
    panic(err)
}

// Get range by rank
players, err := conn.ZRevRange(ctx, "leaderboard", 0, -1)
if err != nil {
    panic(err)
}
fmt.Println(players) // [player2 player3 player1]
```

### Scan Operations

```go
// Scan keys
keys, cursor, err := conn.Scan(ctx, 0, "user:*", 100)
if err != nil {
    panic(err)
}
fmt.Println(keys)

// Scan hash fields
fields, cursor, err := conn.HScan(ctx, "user:123", 0, "*", 100)
if err != nil {
    panic(err)
}
fmt.Println(fields)
```

## Testing

The library includes testing utilities for integration testing:

```go
package main

import (
    "context"
    "testing"

    "github.com/rshelekhov/golib/db/redis"
    "github.com/rshelekhov/golib/db/redis/testutil"
)

func TestRedisOperations(t *testing.T) {
    ctx := context.Background()

    // Create test database
    testDB, err := testutil.NewTestDB(ctx)
    if err != nil {
        t.Fatal(err)
    }
    defer testDB.Close(ctx)

    // Create connection using test database
    conn, err := redis.NewConnection(ctx,
        redis.WithHost(testDB.Host()),
        redis.WithPort(testDB.Port()),
        redis.WithPassword(testDB.Password()),
        redis.WithDB(testDB.DB()),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer conn.Close()

    // Your tests here
    err = conn.Set(ctx, "test", "value", 0)
    if err != nil {
        t.Fatal(err)
    }

    value, err := conn.Get(ctx, "test")
    if err != nil {
        t.Fatal(err)
    }

    if value != "value" {
        t.Errorf("expected 'value', got '%s'", value)
    }
}
```

## Environment Variables for Testing

- `TEST_REDIS_ADDR`: Redis address (e.g., "localhost:6379")
- `TEST_REDIS_PASSWORD`: Redis password
- `TEST_REDIS_DB`: Redis database number

If these are not set, the library will automatically start a Redis container using Docker.

## Architecture

The library follows the same patterns as the MongoDB and PostgreSQL libraries in this repository:

- **Connection Management**: Centralized connection handling with configurable options
- **Interface-Based Design**: Clean separation of concerns through interfaces
- **Transaction Support**: Consistent transaction API across different databases
- **Testing Support**: Built-in testing utilities for integration tests
- **OpenTelemetry Integration**: Observability support out of the box

## Contributing

Please follow the existing patterns in the codebase when contributing.
