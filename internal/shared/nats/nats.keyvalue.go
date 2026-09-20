package natsbroker

import (
	"context"
	"time"

	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsKeyValueStore struct {
	bucket jetstream.KeyValue
}

type INatsKeyValueStore appInterfaces.Cache[string, []byte]

var _ INatsKeyValueStore = (*NatsKeyValueStore)(nil)

func (b Broker) NewKeyValueBucket(ctx context.Context, Name string, durable bool, ttl time.Duration, maxBucketSize int64) (INatsKeyValueStore, error) {
	var storage jetstream.StorageType
	var history int = 1

	//key value is durable already so if user wants pure memory based so only then we need this config
	if !durable {
		storage = jetstream.MemoryStorage
		history = 1
	}
	kv, err := b.js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:   Name,
		Storage:  storage,
		MaxBytes: maxBucketSize,
		History:  uint8(history),
		TTL:      ttl,
	})

	store := NatsKeyValueStore{}
	store.bucket = kv

	return &store, err
}

func (kv NatsKeyValueStore) Delete(ctx context.Context, key string) error {
	return kv.bucket.Delete(ctx, key, nil)
}

func (kv NatsKeyValueStore) Get(ctx context.Context, key string) ([]byte, error) {
	entry, err := kv.bucket.Get(ctx, key)
	return entry.Value(), err
}

func (kv NatsKeyValueStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	_, err := kv.bucket.Put(ctx, key, value)
	return err
}

func (kv NatsKeyValueStore) Update(ctx context.Context, key string, value []byte) error {
	entry, err := kv.bucket.Get(ctx, key)
	if err != nil {
		return err
	}
	_, updationErr := kv.bucket.Update(ctx, key, value, entry.Revision())

	return updationErr
}
