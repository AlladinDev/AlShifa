package internals

import (
	"context"
	"errors"
	"sync"
)

type BusHandler func(ctx context.Context, payload any) (any, error)

type SyncBus struct {
	sc    sync.RWMutex
	Buses map[string]BusHandler
}

func NewSyncBus() *SyncBus {
	return &SyncBus{
		sc:    sync.RWMutex{},
		Buses: map[string]BusHandler{},
	}
}

func (b *SyncBus) RegisterHandler(topic string, handler BusHandler) error {
	b.sc.Lock()
	defer b.sc.Unlock()
	//first check if this topic exists in buses or not if yes throw error as topics should be unique
	_, topicExists := b.Buses[topic]
	if topicExists {
		return errors.New("this topic already exists topics should be unique")
	}

	b.Buses[topic] = handler
	return nil
}

func (b *SyncBus) Request(topic string, ctx context.Context, payload any) (any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	//return nil, errors.New("cancelled request because of done channel  input"
	b.sc.RLock()
	bus, topicExists := b.Buses[topic]
	b.sc.RUnlock()

	if !topicExists {
		return nil, errors.New("this topic doesnt exist")
	}

	res, err := bus(ctx, payload)

	return res, err
}

func (b *SyncBus) HealthCheck() {

}
