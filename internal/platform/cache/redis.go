package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultTimeout = 3 * time.Second

// RedisCache is a thin wrapper around go-redis that satisfies the Cache interface.
// It enforces a small timeout for every operation so runaway requests do not block the app.
type RedisCache struct {
	rdb     *redis.Client
	timeout time.Duration
}

// NewRedis wires a redis client using the supplied connection parameters.
// It performs a health check on startup and exits the program if the connection fails.
func NewRedis(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	return &RedisCache{rdb: client, timeout: defaultTimeout}
}

// Close releases the underlying redis connection pool.
func (c *RedisCache) Close() error { return c.rdb.Close() }

// Get fetches a cached value by key. A missing key returns ("", nil).
func (c *RedisCache) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	v, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return v, err
}

// Set stores a value with the provided TTL.
func (c *RedisCache) Set(key, val string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.rdb.Set(ctx, key, val, ttl).Err()
}

// Del removes one or more keys from the cache.
func (c *RedisCache) Del(keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.rdb.Del(ctx, keys...).Err()
}
