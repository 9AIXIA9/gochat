package kafka

import (
	"gochat/internal/shared/command"
	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func eventToMessage(ev event.Event) *ckafka.Message {
	topic := ev.Topic().String()

	headers := []ckafka.Header{
		{
			Key:   eventIDKey,
			Value: []byte(ev.ID()),
		},
	}

	for key, value := range ev.Headers() {
		headers = append(headers, ckafka.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	return &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Value:     ev.Payload(),
		Timestamp: ev.OccurredAt(),
		Key:       []byte(ev.AggregateID().String()),
		Headers:   headers,
		Opaque:    ev.ID(), // 回调 识别消息
	}
}

func commandToMessage(com command.Command) *ckafka.Message {
	action := com.Action().String()

	headers := []ckafka.Header{
		{
			Key:   commandIDKey,
			Value: []byte(com.ID()),
		},
	}

	for key, value := range com.Headers() {
		headers = append(headers, ckafka.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	return &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &action,
			Partition: ckafka.PartitionAny,
		},
		Value:     com.Payload(),
		Timestamp: com.OccurredAt(),
		Key:       []byte(com.AggregateID().String()),
		Headers:   headers,
		Opaque:    com.ID(), // 回调 识别消息
	}
}
