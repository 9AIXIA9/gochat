package middleware

import (
	"context"
	"gochat/internal/infrastructure/kafka"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	eventIDKey   = "event_id"
	commandIDKey = "command_id"
	receiptIDKey = "receipt_id"
)

func NewDeadLetterErrorMiddlewareWithNamespace(
	creator event.DeadLetterCreator,
	namespace string,
) kafka.ErrorMiddleware {
	return func(next kafka.ErrorHandler) kafka.ErrorHandler {
		return kafka.ErrorHandlerFunc(func(ctx context.Context, err error, message *ckafka.Message) {
			if err := creator.CreateDeadLetter(ctx, convertMessageToEvent(message, namespace), err); err != nil {
				next.Handle(ctx, err, message)
				return
			}
			zap.L().Info(
				"message sent to dead letter queue",
				zap.String("topic", *message.TopicPartition.Topic),
				zap.String("namespace", namespace),
				zap.Int32("partition", message.TopicPartition.Partition),
				zap.Int64("offset", int64(message.TopicPartition.Offset)),
				zap.Binary("key", message.Key),
				zap.Binary("value", message.Value),
			)
		})
	}
}

func convertMessageToEvent(message *ckafka.Message, namespace string) event.Event {
	var id event.ID
	headers := make(map[string]string, len(message.Headers)+2)
	for _, header := range message.Headers {
		headers[header.Key] = string(header.Value)
		if header.Key == eventIDKey || header.Key == commandIDKey || header.Key == receiptIDKey {
			id = event.ID(header.Value)
		}
	}
	if namespace != "" {
		headers["kafka_consumer_group"] = namespace
	}
	if message.TopicPartition.Topic != nil {
		headers["kafka_topic"] = *message.TopicPartition.Topic
	}
	return event.LoadStandardEvent(
		id,
		kernel.ID(message.Key),
		message.Timestamp,
		event.Topic(*message.TopicPartition.Topic),
		message.Value,
		headers,
	)
}
