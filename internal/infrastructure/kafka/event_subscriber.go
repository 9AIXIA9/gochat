package kafka

import (
	"context"
	"fmt"
	"time"

	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var _ event.Subscriber = (*EventSubscriber)(nil)

type EventSubscriber struct {
	consumer    *ckafka.Consumer
	converter   *EventConverter
	handlers    map[event.Topic]event.Handler
	pollTimeout time.Duration
	running     bool
}

// NewEventSubscriber creates a Kafka consumer-based subscriber.
func NewEventSubscriber(commonConfig *CommonConfig, consumerConfig *ConsumerConfig) (*EventSubscriber, error) {
	cCfg, err := getConsumerConfig(commonConfig, consumerConfig)
	if err != nil {
		return nil, err
	}

	consumer, err := ckafka.NewConsumer(cCfg)
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer failed: %w", err)
	}

	return &EventSubscriber{
		consumer:    consumer,
		converter:   &EventConverter{},
		handlers:    make(map[event.Topic]event.Handler),
		pollTimeout: 500 * time.Millisecond,
	}, nil
}

// Subscribe registers a handler for a topic.
func (s *EventSubscriber) Subscribe(topic event.Topic, handler event.Handler) {
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
				e, ok := s.converter.ToEvent(m)
				if !ok {
					// Fallback: build minimal event without payload headers
					e = event.NewStandardEvent(
						"", // unknown id
						kernel.ID(m.Key),
						m.Timestamp,
						event.Topic(*m.TopicPartition.Topic),
						m.Value,
					)
				}

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
					// HandlerFunc error: log and do not commit to allow redelivery
					zap.L().Error("kafka handler error", zap.Error(err), zap.String("topic", topic.String()))
					cancel()
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
