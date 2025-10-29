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

const (
	publisherResultCache = 4096
)

var _ event.Publisher = (*EventPublisher)(nil)

// EventPublisher publishes domain events to Kafka and marks them published upon delivery.
// It relies on an outbox table and only deletes (marks published) after Kafka acks.
type EventPublisher struct {
	publishResultChan chan ckafka.Event
	producer          *ckafka.Producer
	converter         *EventConverter
	publishedMarker   event.PublishedMarker
	deadEventSaver    event.DeadEventSaver // optional

	maxRetries int
	backoff    time.Duration
}

// NewEventPublisher creates a Kafka producer-based publisher.
func NewEventPublisher(cfg *Config, publishedMarker event.PublishedMarker, deadEventSaver event.DeadEventSaver) (*EventPublisher, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid kafka config: %w", err)
	}

	kafkaConfig, err := cfg.toKafkaConfig()
	if err != nil {
		return nil, err
	}

	producer, err := ckafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer failed: %w", err)
	}

	maxRetries := cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}
	backoff := time.Duration(cfg.BackoffMs) * time.Millisecond
	if backoff <= 0 {
		backoff = 200 * time.Millisecond
	}

	return &EventPublisher{
		publishResultChan: make(chan ckafka.Event, publisherResultCache),
		producer:          producer,
		converter:         &EventConverter{},
		publishedMarker:   publishedMarker,
		deadEventSaver:    deadEventSaver,
		maxRetries:        maxRetries,
		backoff:           backoff,
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

func (p *EventPublisher) PublishEvents(events []event.Event) error {
	if len(events) == 0 {
		return nil
	}

	for _, e := range events {
		if err := p.produceWithRetry(p.converter.ToMessage(e)); err != nil {
			zap.L().Error("produce message failed", zap.Error(err))
			// As a last resort, attempt to send to dead letter store
			if p.deadEventSaver != nil {
				_ = p.deadEventSaver.SaveDeadEvent(context.Background(), e, err)
			}
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
				// 失败：从 Opaque 取回原事件，写入死信
				if p.deadEventSaver != nil {
					if e, ok := m.Opaque.(event.Event); ok {
						_ = p.deadEventSaver.SaveDeadEvent(context.Background(), e, m.TopicPartition.Error)
					} else {
						// 兜底：尝试从回执重建
						if e2, ok := p.converter.ToEvent(m); ok {
							_ = p.deadEventSaver.SaveDeadEvent(context.Background(), e2, m.TopicPartition.Error)
						}
					}
				}
				continue
			}

			// 成功：优先用 Opaque 中的事件 ID 标记已发布
			if e, ok := m.Opaque.(event.Event); ok {
				if err := p.publishedMarker.MarkPublished(context.Background(), e.ID()); err != nil {
					zap.L().Error("mark published failed", zap.Error(err), zap.String("event_id", string(e.ID())))
				}
				continue
			}

			// 兜底：尝试从 Headers 读取（大多情况下为空）
			if id, ok := getEventIDFromHeaders(m.Headers); ok {
				if err := p.publishedMarker.MarkPublished(context.Background(), id); err != nil {
					zap.L().Error("mark published failed", zap.Error(err), zap.String("event_id", string(id)))
				}
			} else {
				zap.L().Warn("delivery ok but event_id not available in delivery report")
			}

		case ckafka.Error:
			zap.L().Error("kafka producer error", zap.Error(m))
		default:
			// ignore
		}
	}
}
