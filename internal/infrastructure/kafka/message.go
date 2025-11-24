package kafka

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"strconv"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func getMessage(event event.Event, retry int) *ckafka.Message {
	topic := event.Topic().String()
	m := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          event.Payload(),
		Timestamp:      event.OccurredAt(),
		Key:            []byte(event.AggregateID().String()),
		Headers: []ckafka.Header{
			{
				Key:   "event_id",
				Value: []byte(event.ID()),
			},
			{
				Key:   "retry",
				Value: []byte(strconv.Itoa(retry)),
			},
		},
		Opaque: event.ID(), // 回调 识别消息
	}
	return m
}

func parseMessage(message *ckafka.Message) (event.Event, int) {
	var id event.ID
	var retryTimes int
	foundID, foundRetry := false, false

	for _, header := range message.Headers {
		if header.Key == "event_id" {
			id = event.ID(header.Value)
			foundID = true
			continue
		}
		if header.Key == "retry" {
			times, err := strconv.ParseInt(string(header.Value), 10, 0)
			if err != nil {
				retryTimes = 0
			} else {
				retryTimes = int(times)
			}
			foundRetry = true
			continue
		}
		if foundID && foundRetry {
			break
		}
	}
	return event.LoadStandardEvent(
		id,
		kernel.ID(message.Key),
		message.Timestamp,
		event.Topic(*message.TopicPartition.Topic),
		message.Value,
	), retryTimes
}
