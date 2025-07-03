package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// ConnectionCloser defines the interface for connection management.
type ConnectionCloser interface {
	// Close closes the connection.
	Close() error
	// Client returns the client instance.
	Client() *redis.Client
	// Ping checks the connection to the Redis server.
	Ping(ctx context.Context) error
}

// StringAPI defines the interface for string operations.
type StringAPI interface {
	// Set sets the key to hold the string value.
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	// Get gets the value of key.
	Get(ctx context.Context, key string) (string, error)
	// Del deletes one or more keys.
	Del(ctx context.Context, keys ...string) (int64, error)
	// Exists returns if key exists.
	Exists(ctx context.Context, keys ...string) (int64, error)
	// Expire sets a timeout on key.
	Expire(ctx context.Context, key string, expiration time.Duration) error
	// ExpireAt sets a timeout on key at the given time.
	ExpireAt(ctx context.Context, key string, tm time.Time) error
	// TTL returns the remaining time to live of a key.
	TTL(ctx context.Context, key string) (time.Duration, error)
}

// HashAPI defines the interface for hash operations.
type HashAPI interface {
	// HSet sets field in the hash stored at key to value.
	HSet(ctx context.Context, key string, values ...any) error
	// HGet returns the value associated with field in the hash stored at key.
	HGet(ctx context.Context, key, field string) (string, error)
	// HGetAll returns all fields and values of the hash stored at key.
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	// HDel deletes one or more hash fields.
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	// HExists returns if field is an existing field in the hash stored at key.
	HExists(ctx context.Context, key, field string) (bool, error)
	// HKeys returns all field names in the hash stored at key.
	HKeys(ctx context.Context, key string) ([]string, error)
	// HVals returns all values in the hash stored at key.
	HVals(ctx context.Context, key string) ([]string, error)
	// HLen returns the number of fields in the hash stored at key.
	HLen(ctx context.Context, key string) (int64, error)
}

// ListAPI defines the interface for list operations.
type ListAPI interface {
	// LPush inserts all the specified values at the head of the list stored at key.
	LPush(ctx context.Context, key string, values ...any) (int64, error)
	// RPush inserts all the specified values at the tail of the list stored at key.
	RPush(ctx context.Context, key string, values ...any) (int64, error)
	// LPop removes and returns the first element of the list stored at key.
	LPop(ctx context.Context, key string) (string, error)
	// RPop removes and returns the last element of the list stored at key.
	RPop(ctx context.Context, key string) (string, error)
	// LLen returns the length of the list stored at key.
	LLen(ctx context.Context, key string) (int64, error)
	// LRange returns the specified elements of the list stored at key.
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
}

// SetAPI defines the interface for set operations.
type SetAPI interface {
	// SAdd adds the specified members to the set stored at key.
	SAdd(ctx context.Context, key string, members ...any) (int64, error)
	// SRem removes the specified members from the set stored at key.
	SRem(ctx context.Context, key string, members ...any) (int64, error)
	// SMembers returns all the members of the set value stored at key.
	SMembers(ctx context.Context, key string) ([]string, error)
	// SIsMember returns if member is a member of the set stored at key.
	SIsMember(ctx context.Context, key string, member any) (bool, error)
	// SCard returns the set cardinality (number of elements) of the set stored at key.
	SCard(ctx context.Context, key string) (int64, error)
}

// SortedSetAPI defines the interface for sorted set operations.
type SortedSetAPI interface {
	// ZAdd adds all the specified members with the specified scores to the sorted set stored at key.
	ZAdd(ctx context.Context, key string, members ...redis.Z) (int64, error)
	// ZRem removes the specified members from the sorted set stored at key.
	ZRem(ctx context.Context, key string, members ...any) (int64, error)
	// ZScore returns the score of member in the sorted set at key.
	ZScore(ctx context.Context, key, member string) (float64, error)
	// ZRange returns the specified range of elements in the sorted set stored at key.
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	// ZRevRange returns the specified range of elements in the sorted set stored at key, with scores ordered from high to low.
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	// ZCard returns the sorted set cardinality (number of elements) of the sorted set stored at key.
	ZCard(ctx context.Context, key string) (int64, error)
}

// ScanAPI defines the interface for scan operations.
type ScanAPI interface {
	// Scan iterates the set of keys in the currently selected Redis database.
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
	// HScan iterates fields of Hash types and their associated values.
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error)
	// SScan iterates elements of Set types.
	SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error)
	// ZScan iterates elements of Sorted Set types and their associated scores.
	ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error)
}

// PipelineAPI defines the interface for pipeline operations.
type PipelineAPI interface {
	// Pipeline creates a new pipeline.
	Pipeline() redis.Pipeliner
	// TxPipeline creates a new transaction pipeline.
	TxPipeline() redis.Pipeliner
}

// ConnectionAPI defines the interface for all Redis operations.
type ConnectionAPI interface {
	ConnectionCloser
	StringAPI
	HashAPI
	ListAPI
	SetAPI
	SortedSetAPI
	ScanAPI
	PipelineAPI
}

// TransactionManagerAPI defines the interface for transaction management.
type TransactionManagerAPI interface {
	// GetQueryEngine returns the appropriate query engine based on the context.
	GetQueryEngine(ctx context.Context) QueryEngine
}

// QueryEngine defines the interface for query operations.
type QueryEngine interface {
	StringAPI
	HashAPI
	ListAPI
	SetAPI
	SortedSetAPI
	ScanAPI
}
