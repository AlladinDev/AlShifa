// Package redispack provides utilities for initializing and interacting with Redis.
package redispack

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache is a generic Redis-backed cache
// K = key type (must be comparable)
// V = value type
type RedisCache[K comparable, V any] struct {
	client *redis.Client
	prefix string
}

// NewRedisCache creates a new generic Redis cache instance
func NewRedisCache[K comparable, V any](
	client *redis.Client,
	prefix string,
) *RedisCache[K, V] {
	return &RedisCache[K, V]{
		client: client,
		prefix: prefix,
	}
}

// internal helper to build Redis keys
func (r *RedisCache[K, V]) redisKey(key K) string {
	if r.prefix == "" {
		return fmt.Sprintf("%v", key)
	}
	return fmt.Sprintf("%s:%v", r.prefix, key)
}

// Set stores a value with optional TTL
// ttl <= 0 means no expiration
func (r *RedisCache[K, V]) Set(
	ctx context.Context,
	key K,
	value V,
	ttl time.Duration,
) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(
		ctx,
		r.redisKey(key),
		data,
		ttl,
	).Err()
}

// Get retrieves a value
// bool = cache hit / miss
func (r *RedisCache[K, V]) Get(
	ctx context.Context,
	key K,
) (V, bool, error) {
	var zero V

	data, err := r.client.Get(ctx, r.redisKey(key)).Bytes()
	if err == redis.Nil {
		return zero, false, nil
	}
	if err != nil {
		return zero, false, err
	}

	var value V
	if err := json.Unmarshal(data, &value); err != nil {
		return zero, false, err
	}

	return value, true, nil
}

// Update updates an existing key
// returns error if key does not exist
func (r *RedisCache[K, V]) Update(
	ctx context.Context,
	key K,
	value V,
	ttl time.Duration,
) error {
	exists, err := r.client.Exists(ctx, r.redisKey(key)).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		return fmt.Errorf("key not found")
	}

	return r.Set(ctx, key, value, ttl)
}

// Delete removes a key
func (r *RedisCache[K, V]) Delete(
	ctx context.Context,
	key K,
) error {
	return r.client.Del(ctx, r.redisKey(key)).Err()
}
