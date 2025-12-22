package kafka

import (
	"context"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

const eventIDKey = "event_id"

func WrapEventHandler(eventHandler event.Handler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *ckafka.Message) error {
		return eventHandler.Handle(ctx, toEvent(msg))
	})
}

func toEvent(message *ckafka.Message) event.Event {
	var id event.ID
	headers := make(map[string]string, len(message.Headers)-1)
	for _, header := range message.Headers {
		if header.Key == eventIDKey {
			id = event.ID(header.Value)
			continue
		}
		headers[header.Key] = string(header.Value)
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
