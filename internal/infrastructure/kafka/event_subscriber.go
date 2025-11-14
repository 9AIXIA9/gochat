package kafka

import (
	"context"
	"fmt"
	"time"

	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var _ event.Subscriber = (*EventSubscriber)(nil)

type EventSubscriber struct {
	consumer *ckafka.Consumer
	retrier  *ConsumerRetrier

	handlers    map[event.Topic]event.Handler
	pollTimeout time.Duration
	running     bool
}

// NewEventSubscriber creates a Kafka consumer-based subscriber.
func NewEventSubscriber(
	consumer *ckafka.Consumer,
	retrier *ConsumerRetrier,
) (*EventSubscriber, error) {
	return &EventSubscriber{
		consumer:    consumer,
		retrier:     retrier,
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
					zap.L().Error("kafka handler error", zap.Error(err), zap.String("topic", topic.String()))
					// retry
					if err := s.retrier.Retry(ctx, e, err, retry); err != nil {
						zap.L().Error("retry event failed", zap.Error(err), zap.String("topic", topic.String()))
					}
					cancel()
					// 提交偏移量，防止该消息反复重投造成堵塞
					if _, cErr := s.consumer.CommitMessage(m); cErr != nil {
						zap.L().Warn("commit after dead-letter failed", zap.Error(cErr), zap.String("topic", topic.String()))
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
