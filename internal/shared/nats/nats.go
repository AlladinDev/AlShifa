package natsbroker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	envNATSURL = "NATS_URL"
)

const (
	defaultAckWait       = 30 * time.Second
	defaultMaxDeliver    = 10
	defaultMaxAge        = 30 * 24 * time.Hour
	defaultMaxMsgs       = int64(-1)
	defaultMaxBytes      = int64(-1)
	defaultMaxMsgsPerSub = int64(-1)
	defaultReplicas      = 1
	defaultStorage       = jetstream.FileStorage
)

// -----------------------------------------------------------------------------
// Root Broker
// -----------------------------------------------------------------------------

// Broker owns the application's single NATS connection and JetStream context.
type Broker struct {
	conn *nats.Conn
	js   jetstream.JetStream
}

// Compile-time check for Broker.
var _ appInterfaces.IMessageBroker = (*Broker)(nil)

// New creates the application's NATS connection and JetStream context.
func New() (*Broker, error) {
	url := strings.TrimSpace(os.Getenv(envNATSURL))
	if url == "" {
		url = nats.DefaultURL
	}

	conn, err := nats.Connect(
		url,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second*2),
		nats.Timeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("create JetStream context: %w", err)
	}

	return &Broker{
		conn: conn,
		js:   js,
	}, nil
}

func (b *Broker) JetStream() jetstream.JetStream {
	return b.js
}

// Publish direct publish via JetStream context.
func (b *Broker) Publish(ctx context.Context, message appInterfaces.Message) error {
	if strings.TrimSpace(message.Subject) == "" {
		return errors.New("subject is required")
	}

	msg := nats.NewMsg(message.Subject)
	msg.Data = message.Payload

	for key, value := range message.Headers {
		msg.Header.Set(key, value)
	}

	if message.ID != "" {
		msg.Header.Set("Nats-Msg-Id", message.ID)
	}

	_, err := b.js.PublishMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("broker publish to subject %q: %w", message.Subject, err)
	}

	return nil
}

// Subscribe auto-provisions or binds to a core subject/stream and starts consuming.
func (b *Broker) Subscribe(
	ctx context.Context,
	subject string,
	handler appInterfaces.MessageHandler,
) (appInterfaces.Subscription, error) {
	if strings.TrimSpace(subject) == "" {
		return nil, errors.New("subject is required")
	}

	options := DefaultConsumerOptions("default-" + strings.ToLower(strings.ReplaceAll(subject, ">", "all")))

	// Dynamically bind using stream lookup or implicit subscription
	stream, err := b.js.StreamNameBySubject(ctx, subject)
	if err != nil {
		return nil, fmt.Errorf("resolve stream for subject %q: %w", subject, err)
	}

	sb := &StreamBroker{
		js:         b.js,
		streamName: stream,
	}

	return sb.SubscribeWithOptions(ctx, subject, options, handler)
}

func (b *Broker) Close() error {
	if b.conn != nil {
		b.conn.Close()
	}
	return nil
}

// -----------------------------------------------------------------------------
// Stream Options
// -----------------------------------------------------------------------------

type StreamOptions struct {
	Name       string
	Subjects   []string
	Storage    jetstream.StorageType
	MaxAge     time.Duration
	MaxMsgs    int64
	MaxBytes   int64
	MaxMsgSize int32
	Replicas   int
	Discard    jetstream.DiscardPolicy
}

func DefaultStreamOptions(name string, subjects ...string) StreamOptions {
	return StreamOptions{
		Name:       strings.ToUpper(strings.TrimSpace(name)),
		Subjects:   subjects,
		Storage:    defaultStorage,
		MaxAge:     defaultMaxAge,
		MaxMsgs:    defaultMaxMsgs,
		MaxBytes:   defaultMaxBytes,
		MaxMsgSize: -1,
		Replicas:   defaultReplicas,
		Discard:    jetstream.DiscardOld,
	}
}

// -----------------------------------------------------------------------------
// Stream Broker
// -----------------------------------------------------------------------------

func (b *Broker) ForStream(ctx context.Context, options StreamOptions) (*StreamBroker, error) {
	if strings.TrimSpace(options.Name) == "" {
		return nil, errors.New("stream name is required")
	}

	if len(options.Subjects) == 0 {
		return nil, errors.New("at least one stream subject is required")
	}

	options.Name = strings.ToUpper(strings.TrimSpace(options.Name))

	if options.Storage == 0 {
		options.Storage = jetstream.FileStorage
	}
	if options.MaxAge <= 0 {
		options.MaxAge = defaultMaxAge
	}
	if options.MaxMsgs == 0 {
		options.MaxMsgs = defaultMaxMsgs
	}
	if options.MaxBytes == 0 {
		options.MaxBytes = defaultMaxBytes
	}
	if options.Replicas <= 0 {
		options.Replicas = defaultReplicas
	}

	stream, err := b.js.CreateOrUpdateStream(
		ctx,
		jetstream.StreamConfig{
			Name:       options.Name,
			Subjects:   options.Subjects,
			Storage:    options.Storage,
			MaxAge:     options.MaxAge,
			MaxMsgs:    options.MaxMsgs,
			MaxBytes:   options.MaxBytes,
			MaxMsgSize: options.MaxMsgSize,
			Replicas:   options.Replicas,
			Discard:    options.Discard,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create/update stream %q: %w", options.Name, err)
	}

	return &StreamBroker{
		js:         b.js,
		stream:     stream,
		streamName: options.Name,
	}, nil
}

type StreamBroker struct {
	js         jetstream.JetStream
	stream     jetstream.Stream
	streamName string

	mu            sync.Mutex
	subscriptions []*subscription
}

// Compile-time check for StreamBroker.
var _ appInterfaces.IMessageBroker = (*StreamBroker)(nil)

func (b *StreamBroker) Publish(ctx context.Context, message appInterfaces.Message) error {
	if strings.TrimSpace(message.Subject) == "" {
		return errors.New("subject is required")
	}

	msg := nats.NewMsg(message.Subject)
	msg.Data = message.Payload

	for key, value := range message.Headers {
		msg.Header.Set(key, value)
	}

	if message.ID != "" {
		msg.Header.Set("Nats-Msg-Id", message.ID)
	}

	_, err := b.js.PublishMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("publish to subject %q: %w", message.Subject, err)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Consumer Configuration & Subscription
// -----------------------------------------------------------------------------

type ConsumerOptions struct {
	Durable       string
	AckWait       time.Duration
	MaxDeliver    int
	DeliverPolicy *jetstream.DeliverPolicy
}

func DefaultConsumerOptions(durable string) ConsumerOptions {
	return ConsumerOptions{
		Durable:    durable,
		AckWait:    defaultAckWait,
		MaxDeliver: defaultMaxDeliver,
	}
}

func (b *StreamBroker) Subscribe(
	ctx context.Context,
	subject string,
	handler appInterfaces.MessageHandler,
) (appInterfaces.Subscription, error) {
	sanitizedSubject := strings.ReplaceAll(strings.ToLower(subject), "*", "all")
	sanitizedSubject = strings.ReplaceAll(sanitizedSubject, ">", "all")

	options := DefaultConsumerOptions("default-" + sanitizedSubject)
	return b.SubscribeWithOptions(ctx, subject, options, handler)
}

func (b *StreamBroker) SubscribeWithOptions(
	ctx context.Context,
	subject string,
	options ConsumerOptions,
	handler appInterfaces.MessageHandler,
) (appInterfaces.Subscription, error) {

	if strings.TrimSpace(subject) == "" {
		return nil, errors.New("subject is required")
	}
	if strings.TrimSpace(options.Durable) == "" {
		return nil, errors.New("durable consumer name is required")
	}
	if handler == nil {
		return nil, errors.New("message handler is required")
	}

	if options.AckWait <= 0 {
		options.AckWait = defaultAckWait
	}
	if options.MaxDeliver == 0 {
		options.MaxDeliver = defaultMaxDeliver
	}

	config := jetstream.ConsumerConfig{
		Name:          options.Durable,
		Durable:       options.Durable,
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       options.AckWait,
		MaxDeliver:    options.MaxDeliver,
	}

	if options.DeliverPolicy != nil {
		config.DeliverPolicy = *options.DeliverPolicy
	}

	consumer, err := b.js.CreateOrUpdateConsumer(ctx, b.streamName, config)
	if err != nil {
		return nil, fmt.Errorf("create/update consumer %q on stream %q: %w", options.Durable, b.streamName, err)
	}

	consumeContext, err := consumer.Consume(func(msg jetstream.Msg) {
		delivery := &natsDelivery{msg: msg}

		// Create an isolated context per message processing lifecycle so context
		// cancellation from the caller's Subscribe function won't abruptly interrupt long-running tasks.
		msgCtx, cancel := context.WithTimeout(context.Background(), options.AckWait)
		defer cancel()

		if err := handler(msgCtx, delivery.Message()); err != nil {
			_ = delivery.Nack(msgCtx, 0)
			return
		}

		_ = delivery.Ack(msgCtx)
	})
	if err != nil {
		return nil, fmt.Errorf("start consumer %q: %w", options.Durable, err)
	}

	sub := &subscription{consumeContext: consumeContext}

	b.mu.Lock()
	b.subscriptions = append(b.subscriptions, sub)
	b.mu.Unlock()

	return sub, nil
}

// Close gracefully stops active subscriptions attached to this StreamBroker without dropping the NATS connection.
func (b *StreamBroker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, sub := range b.subscriptions {
		_ = sub.Close()
	}
	b.subscriptions = nil
	return nil
}

// -----------------------------------------------------------------------------
// Delivery & Subscription Implementations
// -----------------------------------------------------------------------------

type natsDelivery struct {
	msg jetstream.Msg
}

// Compile-time check for AcknowledgableDelivery.
var _ appInterfaces.AcknowledgableDelivery = (*natsDelivery)(nil)

func (d *natsDelivery) Message() appInterfaces.Message {
	headers := make(map[string]string)

	for key, values := range d.msg.Headers() {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	messageID := ""
	if header := d.msg.Headers(); header != nil {
		messageID = header.Get("Nats-Msg-Id")
	}

	return appInterfaces.Message{
		ID:          messageID,
		Subject:     d.msg.Subject(),
		Payload:     d.msg.Data(),
		Headers:     headers,
		ContentType: headers["Content-Type"],
	}
}

func (d *natsDelivery) Ack(ctx context.Context) error {
	return d.msg.Ack()
}

func (d *natsDelivery) Nack(ctx context.Context, delay time.Duration) error {
	if delay > 0 {
		return d.msg.NakWithDelay(delay)
	}
	return d.msg.Nak()
}

func (d *natsDelivery) Reject(ctx context.Context) error {
	return d.msg.Term()
}

type subscription struct {
	consumeContext jetstream.ConsumeContext
}

func (s *subscription) Close() error {
	if s == nil || s.consumeContext == nil {
		return nil
	}
	s.consumeContext.Stop()
	return nil
}
