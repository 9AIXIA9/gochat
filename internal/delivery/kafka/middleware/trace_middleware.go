package middleware

import (
	"context"
	"gochat/internal/infrastructure/kafka"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func NewTraceMiddleware(serviceName string) kafka.Middleware {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) (err error) {
			carrier := make(propagation.MapCarrier, len(message.Headers))
			for _, header := range message.Headers {
				carrier[header.Key] = string(header.Value)
			}

			parentCtx := propagator.Extract(ctx, carrier)
			ctxWithTrace, span := tracer.Start(parentCtx, *message.TopicPartition.Topic)
			span.SetAttributes()

			return next.Handle(ctxWithTrace, message)
		})
	}
}
