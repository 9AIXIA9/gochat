package kafka

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type EventConverter struct{}

func getEventIDFromHeaders(headers []ckafka.Header) (event.ID, bool) {
	for _, h := range headers {
		if h.Key == "event_id" {
			return event.ID(h.Value), true
		}
	}
	return "", false
}

func (c *EventConverter) ToMessage(event event.Event) *ckafka.Message {
	topic := event.Topic().String()
	return &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          event.Payload(),
		Timestamp:      event.OccurredAt(),
		Key:            []byte(event.AggregateID().String()),
		Headers: []ckafka.Header{
			{Key: "event_id", Value: []byte(event.ID())},
		},
		Opaque: event, //供回执读取
	}
}

func (c *EventConverter) ToEvent(message *ckafka.Message) (event.Event, bool) {
	id, ok := getEventIDFromHeaders(message.Headers)
	if !ok {
		return nil, false
	}
	return event.NewStandardEvent(
		id,
		kernel.ID(message.Key),
		message.Timestamp,
		event.Topic(*message.TopicPartition.Topic),
		message.Value,
	), true
}
