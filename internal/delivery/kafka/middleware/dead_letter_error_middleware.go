package middleware

import (
	"context"
	"gochat/internal/infrastructure/kafka"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func NewDeadLetterErrorMiddleware(
	creator event.DeadLetterCreator,
) kafka.ErrorMiddleware {
	return func(next kafka.ErrorHandler) kafka.ErrorHandler {
		return kafka.ErrorHandlerFunc(func(ctx context.Context, err error, message *ckafka.Message) {
			if err := creator.CreateDeadLetter(ctx, convertMessageToEvent(message), err); err != nil {
				next.Handle(ctx, err, message)
				return
			}
			zap.L().Info(
				"message sent to dead letter queue",
				zap.String("topic", *message.TopicPartition.Topic),
				zap.Int32("partition", message.TopicPartition.Partition),
				zap.Int64("offset", int64(message.TopicPartition.Offset)),
				zap.ByteString("key", message.Key),
				zap.ByteString("value", message.Value),
			)
		})
	}
}

func convertMessageToEvent(message *ckafka.Message) event.Event {
	var id event.ID
	for _, header := range message.Headers {
		if header.Key == "event_id" {
			id = event.ID(header.Value)
			break
		}
	}
	return event.LoadStandardEvent(
		id,
		kernel.ID(message.Key),
		message.Timestamp,
		event.Topic(*message.TopicPartition.Topic),
		message.Value,
	)
}
