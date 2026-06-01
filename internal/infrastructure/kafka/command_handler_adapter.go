package kafka

import (
	"context"
	"gochat/internal/shared/command"
	"gochat/internal/shared/kernel"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	commandIDKey = "command_id"
)

func WrapCommandHandler(commandHandler command.Handler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *ckafka.Message) error {
		return commandHandler.Handle(ctx, toCommand(msg))
	})
}

func commandToMessage(com command.Command) *ckafka.Message {
	topic := com.Action().String()

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
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Value:     com.Payload(),
		Timestamp: com.OccurredAt(),
		Key:       []byte(com.AggregateID().String()),
		Headers:   headers,
		Opaque:    com.ID(), // 回调 识别消息
	}
}

func toCommand(message *ckafka.Message) command.Command {
	var id command.ID
	headers := make(map[string]string, len(message.Headers)-1)
	for _, header := range message.Headers {
		if header.Key == commandIDKey {
			id = command.ID(header.Value)
			continue
		}
		headers[header.Key] = string(header.Value)
	}

	return command.LoadStandardCommand(
		id,
		kernel.ID(message.Key),
		message.Timestamp,
		command.Action(*message.TopicPartition.Topic),
		message.Value,
		headers,
	)
}
