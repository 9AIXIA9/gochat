package middleware

import (
	"context"
	"time"

	"gochat/internal/infrastructure/kafka"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewTelemetryMiddleware(serviceName string) kafka.Middleware {
	tracer := otel.Tracer(serviceName + "/kafka")
	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) error {
			start := time.Now().UTC()
			ctx, span := tracer.Start(ctx, *message.TopicPartition.Topic, trace.WithSpanKind(trace.SpanKindServer))
			err := next.Handle(ctx, message)

			dur := time.Since(start).Seconds()
			span.SetAttributes(
				attribute.String("kafka.topic", *message.TopicPartition.Topic),
				attribute.Float64("kafka.duration_seconds", dur),
			)
			if err != nil {
				span.RecordError(err)
			}
			span.End()
			return err
		})
	}
}
