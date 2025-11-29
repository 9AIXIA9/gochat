package middleware

import (
	"context"
	"gochat/internal/infrastructure/kafka"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// NewLoggerMiddleware returns a middleware that logs request latency and tracing info.
func NewLoggerMiddleware() kafka.Middleware {
	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) error {
			start := time.Now()

			// Call the next handler in the chain
			err := next.Handle(ctx, message)

			dur := time.Since(start)
			span := trace.SpanFromContext(ctx)
			fields := []zap.Field{zap.Duration("latency", dur)}
			if span != nil && span.SpanContext().IsValid() {
				fields = append(fields,
					zap.String("trace_id", span.SpanContext().TraceID().String()),
					zap.String("span_id", span.SpanContext().SpanID().String()),
				)
			}

			zap.L().Info("kafka event completed", append(fields,
				zap.String("topic", *message.TopicPartition.Topic),
				zap.Bool("error", err != nil),
			)...)

			return err
		})
	}
}
