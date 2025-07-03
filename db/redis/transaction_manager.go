package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type key string

const (
	pipelineKey key = "pipeline"
)

// TransactionManager manages Redis transactions using pipelines.
type TransactionManager struct {
	conn *Connection
}

// Pipeline wraps Redis pipeline to implement QueryEngine interface.
type Pipeline struct {
	pipe redis.Pipeliner
}

// NewTransactionManager creates a new transaction manager.
func NewTransactionManager(conn *Connection) *TransactionManager {
	return &TransactionManager{conn: conn}
}

// GetQueryEngine returns the appropriate query engine based on the context.
// If a pipeline exists in the context, it returns the pipeline.
// Otherwise, it returns the connection.
func (m *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	if pipe, ok := ctx.Value(pipelineKey).(*Pipeline); ok {
		return pipe
	}
	return m.conn
}

// RunTransaction executes the given function within a Redis transaction pipeline.
// Redis transactions are implemented using MULTI/EXEC through pipelines.
func (m *TransactionManager) RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// If it's nested transaction, skip initiating a new one
	if _, ok := ctx.Value(pipelineKey).(*Pipeline); ok {
		return fn(ctx)
	}

	// Create transaction pipeline
	pipe := m.conn.client.TxPipeline()
	pipeline := &Pipeline{pipe: pipe}

	// Set pipeline to context
	ctx = context.WithValue(ctx, pipelineKey, pipeline)

	// Execute the function
	if err := fn(ctx); err != nil {
		// Discard the pipeline on error
		pipe.Discard()
		return fmt.Errorf("transaction execution failed: %w", err)
	}

	// Execute the pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("transaction execution failed: %w", err)
	}

	return nil
}

// RunPipeline executes the given function within a Redis pipeline (non-transactional).
func (m *TransactionManager) RunPipeline(ctx context.Context, fn func(ctx context.Context) error) error {
	// If it's nested pipeline, skip initiating a new one
	if _, ok := ctx.Value(pipelineKey).(*Pipeline); ok {
		return fn(ctx)
	}

	// Create pipeline
	pipe := m.conn.client.Pipeline()
	pipeline := &Pipeline{pipe: pipe}

	// Set pipeline to context
	ctx = context.WithValue(ctx, pipelineKey, pipeline)

	// Execute the function
	if err := fn(ctx); err != nil {
		// Discard the pipeline on error
		pipe.Discard()
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	// Execute the pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}

// Pipeline QueryEngine implementation
func (p *Pipeline) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return p.pipe.Set(ctx, key, value, expiration).Err()
}

func (p *Pipeline) Get(ctx context.Context, key string) (string, error) {
	return p.pipe.Get(ctx, key).Result()
}

func (p *Pipeline) Del(ctx context.Context, keys ...string) (int64, error) {
	return p.pipe.Del(ctx, keys...).Result()
}

func (p *Pipeline) Exists(ctx context.Context, keys ...string) (int64, error) {
	return p.pipe.Exists(ctx, keys...).Result()
}

func (p *Pipeline) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return p.pipe.Expire(ctx, key, expiration).Err()
}

func (p *Pipeline) ExpireAt(ctx context.Context, key string, tm time.Time) error {
	return p.pipe.ExpireAt(ctx, key, tm).Err()
}

func (p *Pipeline) TTL(ctx context.Context, key string) (time.Duration, error) {
	return p.pipe.TTL(ctx, key).Result()
}

func (p *Pipeline) HSet(ctx context.Context, key string, values ...any) error {
	return p.pipe.HSet(ctx, key, values...).Err()
}

func (p *Pipeline) HGet(ctx context.Context, key, field string) (string, error) {
	return p.pipe.HGet(ctx, key, field).Result()
}

func (p *Pipeline) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return p.pipe.HGetAll(ctx, key).Result()
}

func (p *Pipeline) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return p.pipe.HDel(ctx, key, fields...).Result()
}

func (p *Pipeline) HExists(ctx context.Context, key, field string) (bool, error) {
	return p.pipe.HExists(ctx, key, field).Result()
}

func (p *Pipeline) HKeys(ctx context.Context, key string) ([]string, error) {
	return p.pipe.HKeys(ctx, key).Result()
}

func (p *Pipeline) HVals(ctx context.Context, key string) ([]string, error) {
	return p.pipe.HVals(ctx, key).Result()
}

func (p *Pipeline) HLen(ctx context.Context, key string) (int64, error) {
	return p.pipe.HLen(ctx, key).Result()
}

func (p *Pipeline) LPush(ctx context.Context, key string, values ...any) (int64, error) {
	return p.pipe.LPush(ctx, key, values...).Result()
}

func (p *Pipeline) RPush(ctx context.Context, key string, values ...any) (int64, error) {
	return p.pipe.RPush(ctx, key, values...).Result()
}

func (p *Pipeline) LPop(ctx context.Context, key string) (string, error) {
	return p.pipe.LPop(ctx, key).Result()
}

func (p *Pipeline) RPop(ctx context.Context, key string) (string, error) {
	return p.pipe.RPop(ctx, key).Result()
}

func (p *Pipeline) LLen(ctx context.Context, key string) (int64, error) {
	return p.pipe.LLen(ctx, key).Result()
}

func (p *Pipeline) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return p.pipe.LRange(ctx, key, start, stop).Result()
}

func (p *Pipeline) SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return p.pipe.SAdd(ctx, key, members...).Result()
}

func (p *Pipeline) SRem(ctx context.Context, key string, members ...any) (int64, error) {
	return p.pipe.SRem(ctx, key, members...).Result()
}

func (p *Pipeline) SMembers(ctx context.Context, key string) ([]string, error) {
	return p.pipe.SMembers(ctx, key).Result()
}

func (p *Pipeline) SIsMember(ctx context.Context, key string, member any) (bool, error) {
	return p.pipe.SIsMember(ctx, key, member).Result()
}

func (p *Pipeline) SCard(ctx context.Context, key string) (int64, error) {
	return p.pipe.SCard(ctx, key).Result()
}

func (p *Pipeline) ZAdd(ctx context.Context, key string, members ...redis.Z) (int64, error) {
	return p.pipe.ZAdd(ctx, key, members...).Result()
}

func (p *Pipeline) ZRem(ctx context.Context, key string, members ...any) (int64, error) {
	return p.pipe.ZRem(ctx, key, members...).Result()
}

func (p *Pipeline) ZScore(ctx context.Context, key, member string) (float64, error) {
	return p.pipe.ZScore(ctx, key, member).Result()
}

func (p *Pipeline) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return p.pipe.ZRange(ctx, key, start, stop).Result()
}

func (p *Pipeline) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return p.pipe.ZRevRange(ctx, key, start, stop).Result()
}

func (p *Pipeline) ZCard(ctx context.Context, key string) (int64, error) {
	return p.pipe.ZCard(ctx, key).Result()
}

func (p *Pipeline) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return p.pipe.Scan(ctx, cursor, match, count).Result()
}

func (p *Pipeline) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return p.pipe.HScan(ctx, key, cursor, match, count).Result()
}

func (p *Pipeline) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return p.pipe.SScan(ctx, key, cursor, match, count).Result()
}

func (p *Pipeline) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return p.pipe.ZScan(ctx, key, cursor, match, count).Result()
}
