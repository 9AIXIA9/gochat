package middleware

import (
	"context"
	"errors"
	"fmt"
	"gochat/internal/infrastructure/kafka"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ulule/limiter/v3"
	"go.uber.org/zap"
)

var ErrKafkaRateLimited = errors.New("kafka message rate limited")

// NewRateLimitMiddleware applies rate limiting to kafka routes with the shared ulule limiter.
func NewRateLimitMiddleware(limiter *limiter.Limiter, namespace string) kafka.Middleware {
	if limiter == nil {
		return func(next kafka.Handler) kafka.Handler { return next }
	}

	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) error {
			topic := "unknown"
			if message != nil && message.TopicPartition.Topic != nil {
				topic = *message.TopicPartition.Topic
			}
			key := topic
			if namespace != "" {
				key = namespace + ":" + topic
			}

			result, err := limiter.Get(ctx, key)
			if err != nil {
				zap.L().Warn("kafka ratelimit check failed, fail-open", zap.Error(err), zap.String("topic", topic), zap.String("namespace", namespace))
				return next.Handle(ctx, message)
			}
			if result.Reached {
				return fmt.Errorf("%w: namespace=%s topic=%s", ErrKafkaRateLimited, namespace, topic)
			}

			return next.Handle(ctx, message)
		})
	}
}
