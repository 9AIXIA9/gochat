package middleware

import (
	"context"
	"gochat/internal/infrastructure/kafka"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

// NewLoggerMiddleware returns a middleware that logs request latency
func NewLoggerMiddleware() kafka.Middleware {
	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) error {
			start := time.Now().UTC()

			// Call the next handler in the chain
			err := next.Handle(ctx, message)

			dur := time.Since(start)
			zap.L().Info(
				"request completed",
				zap.String("topic", *message.TopicPartition.Topic),
				zap.Duration("latency", dur),
				zap.Error(err),
			)
			return err
		})
	}
}
