package middleware

import (
	"context"
	"gochat/internal/infrastructure/kafka"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

func NewTraceMiddleware(serviceName string) kafka.Middleware {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) (err error) {
			// Build carrier from Kafka headers (handle nil safely)
			carrier := propagation.MapCarrier{}
			if message.Headers != nil {
				for _, header := range message.Headers {
					carrier[header.Key] = string(header.Value)
				}
			}

			// Extract parent context and start a span named by topic
			parentCtx := propagator.Extract(ctx, carrier)
			topic := "unknown"
			if message.TopicPartition.Topic != nil {
				topic = *message.TopicPartition.Topic
			}
			ctxWithTrace, span := tracer.Start(parentCtx, topic)
			// Ensure the span is ended regardless of handler outcome
			defer span.End()

			// Add useful attributes
			span.SetAttributes(
				attribute.String("messaging.system", "kafka"),
				attribute.String("messaging.destination", topic),
				attribute.String("messaging.destination_kind", "topic"),
				attribute.String("messaging.operation", "process"),
				attribute.String("messaging.kafka.partition", message.TopicPartition.String()),
			)

			// Call downstream handler within the span context
			err = next.Handle(ctxWithTrace, message)
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
			}
			return err
		})
	}
}
