package middleware

import (
	"context"
	"time"

	"gochat/internal/infrastructure/kafka"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	inboxKeyPrefix     = "kafka:inbox:"
	inboxEventIDKey    = "event_id"
	inboxProcessingTTL = 60 * time.Second
	inboxProcessedTTL  = 24 * time.Hour
)

// InboxStore abstracts the persistence used to deduplicate Kafka messages.
type InboxStore interface {
	Claim(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, key string) error
	Complete(ctx context.Context, key string, ttl time.Duration) error
}

// NewInboxMiddleware returns a middleware that uses the inbox pattern to implement idempotence.
// It claims an event id before handling. If the claim fails because the key already exists,
// the message is treated as a duplicate and skipped. On success, the key TTL is extended.
// On handler error, the claim is released so the message can be retried.
func NewInboxMiddleware(store InboxStore) kafka.Middleware {
	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) error {
			if store == nil {
				return next.Handle(ctx, message)
			}

			id := getMessageEventID(message)
			if id == "" {
				return next.Handle(ctx, message)
			}

			topic := getMessageTopic(message)
			key := buildInboxKey(topic, id)

			claimed, err := store.Claim(ctx, key, inboxProcessingTTL)
			if err != nil {
				zap.L().Warn("inbox middleware claim failed, proceeding without dedupe",
					zap.String("topic", topic),
					zap.String("event_id", id),
					zap.Error(err),
				)
				return next.Handle(ctx, message)
			}

			if !claimed {
				zap.L().Info("kafka message skipped by inbox middleware (duplicate)",
					zap.String("topic", topic),
					zap.String("event_id", id),
				)
				return nil
			}

			if handlerErr := next.Handle(ctx, message); handlerErr != nil {
				if releaseErr := store.Release(ctx, key); releaseErr != nil {
					zap.L().Warn("inbox middleware failed to release claim after handler error",
						zap.String("topic", topic),
						zap.String("event_id", id),
						zap.Error(releaseErr),
					)
				}
				return handlerErr
			}

			if err := store.Complete(ctx, key, inboxProcessedTTL); err != nil {
				zap.L().Warn("inbox middleware failed to mark message as processed",
					zap.String("topic", topic),
					zap.String("event_id", id),
					zap.Error(err),
				)
			}
			return nil
		})
	}
}

func buildInboxKey(topic, eventID string) string {
	return inboxKeyPrefix + topic + ":" + eventID
}

func getMessageEventID(message *ckafka.Message) string {
	for _, header := range message.Headers {
		if header.Key == inboxEventIDKey {
			return string(header.Value)
		}
	}
	return ""
}

func getMessageTopic(message *ckafka.Message) string {
	if message.TopicPartition.Topic != nil {
		return *message.TopicPartition.Topic
	}
	return "unknown"
}
