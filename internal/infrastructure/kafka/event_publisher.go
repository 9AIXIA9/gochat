package kafka

import (
	"context"
	"fmt"

	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var _ event.Publisher = (*EventPublisher)(nil)

type EventPublisher struct {
	publishResultChan chan ckafka.Event
	producer          *ckafka.Producer
	publishedMarker   event.PublishedMarker
}

func NewEventPublisher(
	config *Config,
	publishedMarker event.PublishedMarker,
) (*EventPublisher, error) {
	producer, err := ckafka.NewProducer(convertToMap(config))
	if err != nil {
		return nil, fmt.Errorf("create kafka producer failed: %w", err)
	}

	return &EventPublisher{
		publishResultChan: make(chan ckafka.Event, 512),
		producer:          producer,
		publishedMarker:   publishedMarker,
	}, nil
}

func (p *EventPublisher) Publish(event event.Event) error {
	if event == nil {
		return nil
	}

	message := getMessage(event, 0)
	if err := p.producer.Produce(message, p.publishResultChan); err != nil {
		zap.L().Error("produce message failed", zap.Error(err))
		return err
	}
	return nil
}

func (p *EventPublisher) Publishes(events []event.Event) error {
	if len(events) == 0 {
		return nil
	}

	for _, e := range events {
		if err := p.Publish(e); err != nil {
			return err
		}
	}
	return nil
}

func (p *EventPublisher) processSendingResponse() {
	for result := range p.publishResultChan {
		switch message := result.(type) {
		case *ckafka.Message:
			if message.TopicPartition.Error != nil {
				continue
			}

			id := message.Opaque.(event.ID)
			if err := p.publishedMarker.MarkPublished(context.Background(), id); err != nil {
				zap.L().Error(
					"mark event published failed",
					zap.Error(err),
					zap.String("event_id", id.String()),
				)
			}
		case ckafka.Error:
			zap.L().Error("kafka producer error", zap.Error(message))
		default:
			// ignore
		}
	}
}

func (p *EventPublisher) Start() {
	go p.processSendingResponse()
}

func (p *EventPublisher) Close() {
	p.producer.Close()
	close(p.publishResultChan)
}
