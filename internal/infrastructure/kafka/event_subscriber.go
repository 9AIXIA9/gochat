package kafka

import (
	"context"
	"fmt"
	"time"

	"gochat/internal/infrastructure/prometheus"
	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

//TODO 添加中间件而不是直接在handle里处理链路和指标

const (
	maxRetries = 3
)

var _ event.Subscriber = (*EventSubscriber)(nil)

type EventSubscriber struct {
	ctx    context.Context
	cancel func()

	consumer        *ckafka.Consumer
	producer        *ckafka.Producer
	deadLetterSaver event.DeadLetterSaver

	handlers    map[event.Topic]event.Handler
	pollTimeout time.Duration
	running     bool
	metrics     *prometheus.Metrics
}

// NewEventSubscriber creates a Kafka consumer-based subscriber.
func NewEventSubscriber(config *Config, saver event.DeadLetterSaver, metrics *prometheus.Metrics) (*EventSubscriber, error) {
	consumer, err := ckafka.NewConsumer(getConsumerConfigMap(config))
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer failed: %w", err)
	}

	producer, err := ckafka.NewProducer(getProducerConfigMap(config))
	if err != nil {
		return nil, fmt.Errorf("create kafka producer failed: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &EventSubscriber{
		ctx:             ctx,
		cancel:          cancel,
		consumer:        consumer,
		producer:        producer,
		deadLetterSaver: saver,
		handlers:        make(map[event.Topic]event.Handler),
		pollTimeout:     500 * time.Millisecond,
		running:         false,
		metrics:         metrics,
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
func (s *EventSubscriber) Start(serviceName string) error {
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

	tracer := otel.Tracer(serviceName + "/kafka-subscriber")

	go func() {
		for s.running {
			// Handle ctx cancellation
			select {
			case <-s.ctx.Done():
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
				start := time.Now()
				ctx, span := tracer.Start(context.Background(), topic.String(), trace.WithSpanKind(trace.SpanKindConsumer))
				timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

				if err := handler.Handle(timeoutCtx, e); err != nil {
					zap.L().Error(
						"kafka handler error",
						zap.Error(err),
						zap.String("id", e.ID().String()),
						zap.String("topic", topic.String()),
						zap.Int("retry", retry),
					)
					span.RecordError(err)

					if s.metrics != nil {
						s.metrics.KafkaConsumed.WithLabelValues(topic.String(), "error").Inc()
						s.metrics.KafkaHandleDur.WithLabelValues(topic.String()).Observe(time.Since(start).Seconds())
					}

					cancel()
					span.SetAttributes(
						attribute.String("kafka.topic", topic.String()),
						attribute.String("event.id", e.ID().String()),
						attribute.Int("event.retry", retry),
						attribute.Float64("event.handle_duration_ms", float64(time.Since(start).Milliseconds())),
					)
					span.End()

					// 提交偏移量，防止该消息反复重投造成堵塞
					if _, cErr := s.consumer.CommitMessage(m); cErr != nil {
						zap.L().Warn(
							"commit failed",
							zap.Error(cErr),
							zap.String("topic", topic.String()),
						)
					}
					s.handleFailedEvent(e, err, retry)
					continue
				}
				cancel()
				span.SetAttributes(
					attribute.String("kafka.topic", topic.String()),
					attribute.String("event.id", e.ID().String()),
					attribute.Int("event.retry", retry),
					attribute.Float64("event.handle_duration_ms", float64(time.Since(start).Milliseconds())),
				)
				span.End()

				if s.metrics != nil {
					s.metrics.KafkaConsumed.WithLabelValues(topic.String(), "success").Inc()
					s.metrics.KafkaHandleDur.WithLabelValues(topic.String()).Observe(time.Since(start).Seconds())
				}

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

func (s *EventSubscriber) handleFailedEvent(ev event.Event, reason error, retry int) {
	if retry >= maxRetries {
		if err := s.deadLetterSaver.SaveDeadLetter(context.Background(), ev, reason); err != nil {
			zap.L().Error(
				"retry >= maxRetries and save dead letter failed",
				zap.String("id", ev.ID().String()),
				zap.String("topic", ev.Topic().String()),
				zap.Int("retry", retry),
				zap.ByteString("payload", ev.Payload()),
				zap.Error(err),
			)
			return
		}
		zap.L().Debug("failed event's retries >= maxRetires, it was saved to dead letter")
		return
	}

	if err := s.producer.Produce(getMessage(ev, retry+1), nil); err != nil {
		if err := s.deadLetterSaver.SaveDeadLetter(context.Background(), ev, reason); err != nil {
			zap.L().Error(
				"republish failed and save dead letter failed",
				zap.String("id", ev.ID().String()),
				zap.String("topic", ev.Topic().String()),
				zap.Int("retry", retry),
				zap.ByteString("payload", ev.Payload()),
				zap.Error(err),
			)
			return
		}
		zap.L().Debug("republish failed, it was saved to dead letter")
		return
	}
}

func (s *EventSubscriber) Close() {
	s.running = false
	s.cancel()
	if s.consumer != nil {
		_ = s.consumer.Close()
	}
}
