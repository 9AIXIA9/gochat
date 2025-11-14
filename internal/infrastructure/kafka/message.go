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
		Opaque:         event, //供回执读取
	}
	m.Headers = append(m.Headers,
		ckafka.Header{
			Key:   "event_id",
			Value: []byte(event.ID()),
		},
		ckafka.Header{
			Key:   "retry_times",
			Value: []byte(strconv.Itoa(retry)),
		},
	)
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
		}
		if header.Key == "retry_times" {
			times, err := strconv.ParseInt(string(header.Value), 10, 0)
			if err != nil {
				retryTimes = 0
			} else {
				retryTimes = int(times)
			}
			foundRetry = true
		}
		if foundID && foundRetry {
			break
		}
	}
	return event.NewStandardEvent(
		id,
		kernel.ID(message.Key),
		message.Timestamp,
		event.Topic(*message.TopicPartition.Topic),
		message.Value,
	), retryTimes
}
