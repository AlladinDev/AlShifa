package interfaces

import (
	"context"
	"time"
)

type Message struct {
	ID          string
	Subject     string
	Key         []byte
	Payload     []byte
	Headers     map[string]string
	ContentType string
}

type MessageHandler func(ctx context.Context, message Message) error

type Subscription interface {
	Close() error
}

type IMessageBroker interface {
	Publish(ctx context.Context, message Message) error

	Subscribe(
		ctx context.Context,
		subject string,
		handler MessageHandler,
	) (Subscription, error)

	Close() error
}

// Optional interface for brokers that support explicit acknowledgments.
type AcknowledgableDelivery interface {
	Message() Message
	Ack(ctx context.Context) error
	Nack(ctx context.Context, delay time.Duration) error
	Reject(ctx context.Context) error
}
