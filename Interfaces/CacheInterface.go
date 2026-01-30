// Package interfaces provides abstractions for caching mechanisms.
package interfaces

import (
	"context"
	"time"
)

type ICache[K comparable, V any] interface {
	Set(ctx context.Context, key K, value V, ttl time.Duration) error
	Get(ctx context.Context, key K) (V, bool, error)
	Update(ctx context.Context, key K, value V, ttl time.Duration) error
	Delete(ctx context.Context, key K) error
}
