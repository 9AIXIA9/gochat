package middleware

import (
	"context"
	"time"

	"gochat/internal/infrastructure/kafka"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	goRedis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	inboxKeyPrefix  = "kafka:inbox:"
	inboxEventIDKey = "event_id"
)

// NewInboxMiddleware returns a middleware that uses the inbox pattern to implement idempotence.
// It reserves an event id in Redis with SetNX before handling. If reservation fails (already
// exists) the message is considered duplicate and skipped. On success the middleware extends
// the key TTL after handler success; on handler error the reservation is removed to allow retry.
func NewInboxMiddleware(client *goRedis.Client) kafka.Middleware {
	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) error {
			var id string
			for _, h := range message.Headers {
				if h.Key == inboxEventIDKey {
					id = string(h.Value)
					break
				}
			}

			// If no event id provided, cannot dedupe — continue normally.
			if id == "" || client == nil {
				return next.Handle(ctx, message)
			}

			topic := "unknown"
			if message.TopicPartition.Topic != nil {
				topic = *message.TopicPartition.Topic
			}

			key := inboxKeyPrefix + topic + ":" + id

			// Try to reserve the event id for a short processing window.
			reserved, err := client.SetNX(ctx, key, "processing", 60*time.Second).Result()
			if err != nil {
				// Redis unavailable — log and continue to avoid blocking consumption.
				zap.L().Warn("inbox middleware redis error, proceeding without dedupe", zap.Error(err))
				return next.Handle(ctx, message)
			}

			if !reserved {
				// Already processing or processed — skip handling to achieve idempotence.
				zap.L().Info("kafka message skipped by inbox middleware (duplicate)",
					zap.String("topic", topic),
					zap.String("event_id", id),
				)
				return nil
			}

			// We reserved the id — handle the message.
			handlerErr := next.Handle(ctx, message)
			if handlerErr != nil {
				// Remove reservation so the message can be retried later.
				if derr := client.Del(ctx, key).Err(); derr != nil {
					zap.L().Warn("inbox middleware failed to delete redis key after handler error", zap.Error(derr))
				}
				return handlerErr
			}

			// On success, keep a processed marker for a longer time to prevent reprocessing.
			if err := client.Expire(ctx, key, 24*time.Hour).Err(); err != nil {
				zap.L().Warn("inbox middleware failed to set ttl on processed key", zap.Error(err))
			}
			return nil
		})
	}
}
