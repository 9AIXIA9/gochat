package kafka

import (
	"context"
	"fmt"
	"time"

	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

//TODO 重试投递到死信队列

var _ event.Subscriber = (*EventSubscriber)(nil)

type EventSubscriber struct {
	consumer *ckafka.Consumer

	handlers    map[event.Topic]event.Handler
	pollTimeout time.Duration
	running     bool
}

// NewEventSubscriber creates a Kafka consumer-based subscriber.
func NewEventSubscriber(config *Config) (*EventSubscriber, error) {
	consumer, err := ckafka.NewConsumer(convertToMap(config))
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer failed: %w", err)
	}

	return &EventSubscriber{
		consumer:    consumer,
		handlers:    make(map[event.Topic]event.Handler),
		pollTimeout: 500 * time.Millisecond,
		running:     false,
	}, nil
}

// Subscribe registers a handler for a topic.
func (s *EventSubscriber) Subscribe(topic event.Topic, handler event.Handler) {
	if _, ok := s.handlers[topic]; ok {
		zap.L().Warn("handler for topic already exists, overwriting", zap.String("topic", topic.String()))
	}
	s.handlers[topic] = handler
}

// Start begins polling and dispatching messages.
func (s *EventSubscriber) Start(ctx context.Context) error {
	if s.running {
		return nil
	}

	// Build topic list
	var topics []string
	for t := range s.handlers {
		topics = append(topics, t.String())
	}
	if len(topics) == 0 {
		return fmt.Errorf("no topics subscribed")
	}

	if err := s.consumer.SubscribeTopics(topics, nil); err != nil {
		return fmt.Errorf("subscribe topics failed: %w", err)
	}

	s.running = true

	go func() {
		for s.running {
			// Handle ctx cancellation
			select {
			case <-ctx.Done():
				s.running = false
				return
			default:
			}

			ev := s.consumer.Poll(int(s.pollTimeout.Milliseconds()))
			if ev == nil {
				continue
			}

			switch m := ev.(type) {
			case *ckafka.Message:
				// Convert message to event
				e, retry := parseMessage(m)

				topic := e.Topic()
				handler, exists := s.handlers[topic]
				if !exists {
					// No handler: commit and continue
					if _, err := s.consumer.CommitMessage(m); err != nil {
						zap.L().Warn("commit message without handler failed", zap.Error(err))
					}
					continue
				}

				// Dispatch with a per-message context (inherits parent)
				msgCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				if err := handler.Handle(msgCtx, e); err != nil {
					zap.L().Error(
						"kafka handler error",
						zap.Error(err),
						zap.String("id", e.ID().String()),
						zap.String("topic", topic.String()),
						zap.Int("retry", retry),
						zap.ByteString("payload", e.Payload()),
					)
					cancel()
					// 提交偏移量，防止该消息反复重投造成堵塞
					if _, cErr := s.consumer.CommitMessage(m); cErr != nil {
						zap.L().Warn(
							"commit failed",
							zap.Error(cErr),
							zap.String("topic", topic.String()),
						)
					}
					continue
				}
				cancel()

				// Success: manual commit
				if _, err := s.consumer.CommitMessage(m); err != nil {
					zap.L().Warn("commit message failed", zap.Error(err))
				}

			case ckafka.Error:
				// Consumer-level error
				zap.L().Error("kafka consumer error", zap.Error(m))
			default:
				// ignore other events
			}
		}
	}()
	return nil
}

func (s *EventSubscriber) Close() {
	s.running = false
	if s.consumer != nil {
		_ = s.consumer.Close()
	}
}
