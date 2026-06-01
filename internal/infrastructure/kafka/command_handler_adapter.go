package kafka

import (
	"context"
	"gochat/internal/shared/command"
	"gochat/internal/shared/kernel"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	commandIDKey = "command_id"
	receiptIDKey = "receipt_id"
)

func WrapCommandHandler(commandHandler command.Handler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *ckafka.Message) error {
		return commandHandler.Handle(ctx, toCommand(msg))
	})
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
