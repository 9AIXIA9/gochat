package kafka

import (
	"context"
	"fmt"
	"strings"

	"gochat/internal/shared/command"
	"gochat/internal/shared/kernel"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	receiptIDKey          = "receipt_id"
	receiptTopicSeparator = "."
)

func WrapReceiptHandler(receiptHandler command.ReceiptHandler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *ckafka.Message) error {
		receipt, err := toReceipt(msg)
		if err != nil {
			return err
		}
		return receiptHandler.Handle(ctx, receipt)
	})
}

func receiptToMessage(receipt command.Receipt) *ckafka.Message {
	topic := receipt.Action().String() + "." + receipt.Status().String()
	headers := []ckafka.Header{
		{
			Key:   receiptIDKey,
			Value: []byte(receipt.ID()),
		},
		{
			Key:   commandIDKey,
			Value: []byte(receipt.CommandID()),
		},
	}

	for key, value := range receipt.Headers() {
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
		Value:     nil,
		Timestamp: receipt.OccurredAt(),
		Key:       []byte(receipt.AggregateID().String()),
		Headers:   headers,
		Opaque:    receipt.ID(), // 回调 识别消息
	}
}

func toReceipt(message *ckafka.Message) (command.Receipt, error) {
	var id command.ReceiptID
	headers := make(map[string]string, len(message.Headers))
	for _, header := range message.Headers {
		if header.Key == receiptIDKey {
			id = command.ReceiptID(header.Value)
			continue
		}
		headers[header.Key] = string(header.Value)
	}

	action, status, err := mustParseReceiptTopic(message.TopicPartition.Topic)
	if err != nil {
		return nil, err
	}
	return command.LoadStandardReceipt(
		id,
		command.ID(headers[commandIDKey]),
		kernel.ID(message.Key),
		message.Timestamp,
		action,
		status,
		headers,
	), nil
}

func parseReceiptTopic(topic *string) (command.Action, command.ReceiptStatus) {
	if topic == nil || *topic == "" {
		return "", ""
	}
	idx := strings.LastIndex(*topic, receiptTopicSeparator)
	if idx <= 0 || idx >= len(*topic)-1 {
		return command.Action(*topic), ""
	}
	return command.Action((*topic)[:idx]), command.ReceiptStatus((*topic)[idx+1:])
}

func mustParseReceiptTopic(topic *string) (command.Action, command.ReceiptStatus, error) {
	action, status := parseReceiptTopic(topic)
	if action == "" || status == "" {
		return "", "", fmt.Errorf("invalid receipt topic")
	}
	return action, status, nil
}
