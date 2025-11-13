package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var _ event.Publisher = (*EventPublisher)(nil)

// EventPublisher publishes domain events to Kafka and marks them published upon delivery.
// It relies on an outbox table and only deletes (marks published) after Kafka acks.
type EventPublisher struct {
	publishResultChan chan ckafka.Event
	producer          *ckafka.Producer
	converter         *EventConverter
	publishedMarker   event.PublishedMarker

	maxRetries int
	backoff    time.Duration
}

// NewEventPublisher creates a Kafka producer-based publisher.
func NewEventPublisher(
	commonConfig *CommonConfig,
	producerConfig *ProducerConfig,
	publishedMarker event.PublishedMarker,
) (*EventPublisher, error) {
	kafkaConfig, err := getProducerConfig(commonConfig, producerConfig)
	if err != nil {
		return nil, err
	}

	producer, err := ckafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer failed: %w", err)
	}

	return &EventPublisher{
		publishResultChan: make(chan ckafka.Event, 4096),
		producer:          producer,
		converter:         &EventConverter{},
		publishedMarker:   publishedMarker,
		maxRetries:        producerConfig.MaxRetries,
		backoff:           time.Duration(producerConfig.BackoffMs) * time.Millisecond,
	}, nil
}

// Start background delivery handling.
func (p *EventPublisher) Start() {
	go p.processSendingResponse()
}

// Close flushes and closes the producer and delivery channel.
func (p *EventPublisher) Close() {
	_ = p.producer.Flush(int((3 * time.Second).Milliseconds()))
	p.producer.Close()
	// Producer won't write to our private delivery channel after Close/Flush
	close(p.publishResultChan)
}

func (p *EventPublisher) Publish(event event.Event) error {
	if event == nil {
		return nil
	}

	if err := p.produceWithRetry(p.converter.ToMessage(event)); err != nil {
		zap.L().Error("produce message failed", zap.Error(err))
	}
	return nil
}

func (p *EventPublisher) Publishes(events []event.Event) error {
	if len(events) == 0 {
		return nil
	}

	for _, e := range events {
		if err := p.produceWithRetry(p.converter.ToMessage(e)); err != nil {
			zap.L().Error("produce message failed", zap.Error(err))
		}
	}
	return nil
}

func (p *EventPublisher) produceWithRetry(msg *ckafka.Message) error {
	var lastErr error
	for attempt := 0; attempt <= p.maxRetries; attempt++ {
		lastErr = p.producer.Produce(msg, p.publishResultChan)
		if lastErr == nil {
			return nil
		}
		var kerr ckafka.Error
		if errors.As(lastErr, &kerr) && kerr.Code() == ckafka.ErrQueueFull {
			// back off and let delivery goroutine drain
			time.Sleep(p.backoff)
			continue
		}
		// non-retryable error
		break
	}
	return lastErr
}

func (p *EventPublisher) processSendingResponse() {
	for ev := range p.publishResultChan {
		switch m := ev.(type) {
		case *ckafka.Message:
			if m.TopicPartition.Error != nil {
				continue
			}

			// 成功：优先用 Opaque 中的事件 MessageID 标记已发布
			if e, ok := m.Opaque.(event.Event); ok {
				if err := p.publishedMarker.MarkPublished(context.Background(), e.ID()); err != nil {
					zap.L().Error("mark published failed", zap.Error(err), zap.String("event_id", e.ID().String()))
				}
				continue
			}
		case ckafka.Error:
			zap.L().Error("kafka producer error", zap.Error(m))
		default:
			// ignore
		}
	}
}
