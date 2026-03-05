package middleware

import (
	"context"
	"gochat/internal/infrastructure/kafka"
	"gochat/pkg/ctxutil"
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

			fields := []zap.Field{
				zap.String("topic", *message.TopicPartition.Topic),
				zap.Duration("latency", dur),
			}

			if traceID, spanID := ctxutil.SpanIDAndTraceIDFrom(ctx); traceID != "" && spanID != "" {
				fields = append(fields,
					zap.String("trace_id", traceID),
					zap.String("span_id", spanID),
				)
			}

			if err != nil {
				fields = append(fields,
					zap.Error(err),
				)
			}

			zap.L().Info(
				"request completed",
				fields...,
			)
			return err
		})
	}
}
