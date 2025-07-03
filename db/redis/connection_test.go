package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rshelekhov/golib/db/redis/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisConnection(t *testing.T) {
	ctx := context.Background()

	// Create test database
	testDB, err := testutil.NewTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	// Create connection using test database
	conn, err := NewConnection(ctx,
		WithHost(testDB.Host()),
		WithPort(testDB.Port()),
		WithPassword(testDB.Password()),
		WithDB(testDB.DB()),
		WithTracing(false), // Disable tracing for tests
	)
	require.NoError(t, err)
	defer conn.Close()

	t.Run("String operations", func(t *testing.T) {
		// Test Set/Get
		err := conn.Set(ctx, "test_key", "test_value", time.Hour)
		require.NoError(t, err)

		value, err := conn.Get(ctx, "test_key")
		require.NoError(t, err)
		assert.Equal(t, "test_value", value)

		// Test Exists
		exists, err := conn.Exists(ctx, "test_key")
		require.NoError(t, err)
		assert.Equal(t, int64(1), exists)

		// Test TTL
		ttl, err := conn.TTL(ctx, "test_key")
		require.NoError(t, err)
		assert.True(t, ttl > 0)

		// Test Delete
		deleted, err := conn.Del(ctx, "test_key")
		require.NoError(t, err)
		assert.Equal(t, int64(1), deleted)

		// Test key doesn't exist
		exists, err = conn.Exists(ctx, "test_key")
		require.NoError(t, err)
		assert.Equal(t, int64(0), exists)
	})

	t.Run("Hash operations", func(t *testing.T) {
		// Test HSet
		err := conn.HSet(ctx, "user:123", "name", "John", "age", "30")
		require.NoError(t, err)

		// Test HGet
		name, err := conn.HGet(ctx, "user:123", "name")
		require.NoError(t, err)
		assert.Equal(t, "John", name)

		// Test HGetAll
		user, err := conn.HGetAll(ctx, "user:123")
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"name": "John", "age": "30"}, user)

		// Test HExists
		exists, err := conn.HExists(ctx, "user:123", "name")
		require.NoError(t, err)
		assert.True(t, exists)

		// Test HKeys
		keys, err := conn.HKeys(ctx, "user:123")
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"name", "age"}, keys)

		// Test HLen
		length, err := conn.HLen(ctx, "user:123")
		require.NoError(t, err)
		assert.Equal(t, int64(2), length)

		// Test HDel
		deleted, err := conn.HDel(ctx, "user:123", "age")
		require.NoError(t, err)
		assert.Equal(t, int64(1), deleted)

		// Cleanup
		_, err = conn.Del(ctx, "user:123")
		require.NoError(t, err)
	})

	t.Run("List operations", func(t *testing.T) {
		// Test LPush
		count, err := conn.LPush(ctx, "mylist", "item1", "item2", "item3")
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// Test LLen
		length, err := conn.LLen(ctx, "mylist")
		require.NoError(t, err)
		assert.Equal(t, int64(3), length)

		// Test LRange
		items, err := conn.LRange(ctx, "mylist", 0, -1)
		require.NoError(t, err)
		assert.Equal(t, []string{"item3", "item2", "item1"}, items)

		// Test LPop
		item, err := conn.LPop(ctx, "mylist")
		require.NoError(t, err)
		assert.Equal(t, "item3", item)

		// Test RPush
		count, err = conn.RPush(ctx, "mylist", "item4")
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// Test RPop
		item, err = conn.RPop(ctx, "mylist")
		require.NoError(t, err)
		assert.Equal(t, "item4", item)

		// Cleanup
		_, err = conn.Del(ctx, "mylist")
		require.NoError(t, err)
	})

	t.Run("Set operations", func(t *testing.T) {
		// Test SAdd
		count, err := conn.SAdd(ctx, "myset", "member1", "member2", "member3")
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// Test SCard
		card, err := conn.SCard(ctx, "myset")
		require.NoError(t, err)
		assert.Equal(t, int64(3), card)

		// Test SIsMember
		isMember, err := conn.SIsMember(ctx, "myset", "member1")
		require.NoError(t, err)
		assert.True(t, isMember)

		// Test SMembers
		members, err := conn.SMembers(ctx, "myset")
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"member1", "member2", "member3"}, members)

		// Test SRem
		removed, err := conn.SRem(ctx, "myset", "member1")
		require.NoError(t, err)
		assert.Equal(t, int64(1), removed)

		// Cleanup
		_, err = conn.Del(ctx, "myset")
		require.NoError(t, err)
	})

	t.Run("Sorted Set operations", func(t *testing.T) {
		// Test ZAdd
		count, err := conn.ZAdd(ctx, "leaderboard",
			redis.Z{Score: 100, Member: "player1"},
			redis.Z{Score: 200, Member: "player2"},
			redis.Z{Score: 150, Member: "player3"},
		)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// Test ZCard
		card, err := conn.ZCard(ctx, "leaderboard")
		require.NoError(t, err)
		assert.Equal(t, int64(3), card)

		// Test ZScore
		score, err := conn.ZScore(ctx, "leaderboard", "player1")
		require.NoError(t, err)
		assert.Equal(t, float64(100), score)

		// Test ZRange
		players, err := conn.ZRange(ctx, "leaderboard", 0, -1)
		require.NoError(t, err)
		assert.Equal(t, []string{"player1", "player3", "player2"}, players)

		// Test ZRevRange
		players, err = conn.ZRevRange(ctx, "leaderboard", 0, -1)
		require.NoError(t, err)
		assert.Equal(t, []string{"player2", "player3", "player1"}, players)

		// Test ZRem
		removed, err := conn.ZRem(ctx, "leaderboard", "player1")
		require.NoError(t, err)
		assert.Equal(t, int64(1), removed)

		// Cleanup
		_, err = conn.Del(ctx, "leaderboard")
		require.NoError(t, err)
	})
}

func TestTransactionManager(t *testing.T) {
	ctx := context.Background()

	// Create test database
	testDB, err := testutil.NewTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	// Create connection using test database
	connAPI, err := NewConnection(ctx,
		WithHost(testDB.Host()),
		WithPort(testDB.Port()),
		WithPassword(testDB.Password()),
		WithDB(testDB.DB()),
		WithTracing(false), // Disable tracing for tests
	)
	require.NoError(t, err)
	defer connAPI.Close()

	conn := connAPI.(*Connection)
	tm := NewTransactionManager(conn)

	t.Run("Transaction operations", func(t *testing.T) {
		err := tm.RunTransaction(ctx, func(ctx context.Context) error {
			engine := tm.GetQueryEngine(ctx)

			// All operations will be queued in the transaction
			if err := engine.Set(ctx, "key1", "value1", 0); err != nil {
				return err
			}

			if err := engine.Set(ctx, "key2", "value2", 0); err != nil {
				return err
			}

			return nil
		})
		require.NoError(t, err)

		// Check that both keys were set
		value1, err := conn.Get(ctx, "key1")
		require.NoError(t, err)
		assert.Equal(t, "value1", value1)

		value2, err := conn.Get(ctx, "key2")
		require.NoError(t, err)
		assert.Equal(t, "value2", value2)

		// Cleanup
		_, err = conn.Del(ctx, "key1", "key2")
		require.NoError(t, err)
	})

	t.Run("Pipeline operations", func(t *testing.T) {
		err := tm.RunPipeline(ctx, func(ctx context.Context) error {
			engine := tm.GetQueryEngine(ctx)

			// All operations will be batched
			if err := engine.Set(ctx, "key3", "value3", 0); err != nil {
				return err
			}

			if err := engine.Set(ctx, "key4", "value4", 0); err != nil {
				return err
			}

			return nil
		})
		require.NoError(t, err)

		// Check that both keys were set
		value3, err := conn.Get(ctx, "key3")
		require.NoError(t, err)
		assert.Equal(t, "value3", value3)

		value4, err := conn.Get(ctx, "key4")
		require.NoError(t, err)
		assert.Equal(t, "value4", value4)

		// Cleanup
		_, err = conn.Del(ctx, "key3", "key4")
		require.NoError(t, err)
	})
}
