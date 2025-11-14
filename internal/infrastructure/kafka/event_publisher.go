package kafka

import (
	"context"

	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var _ event.Publisher = (*EventPublisher)(nil)

type EventPublisher struct {
	publishResultChan chan ckafka.Event
	producer          *ProducerWithRetry
	publishedMarker   event.PublishedMarker
}

func NewEventPublisher(
	producer *ProducerWithRetry,
	publishedMarker event.PublishedMarker,
) (*EventPublisher, error) {
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
	if err := p.producer.produceWithRetry(message, p.publishResultChan); err != nil {
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
			ev, _ := parseMessage(message)
			if err := p.publishedMarker.MarkPublished(context.Background(), ev.ID()); err != nil {
				zap.L().Error(
					"mark event published failed",
					zap.Error(err),
					zap.String("event_id", ev.ID().String()),
				)
			}
		case ckafka.Error:
			zap.L().Error("kafka producer error", zap.Error(message))
		default:
			// ignore
		}
	}
}

// Start background delivery handling.
func (p *EventPublisher) Start() {
	go p.processSendingResponse()
}

// Close flushes and closes the producer and delivery channel.
func (p *EventPublisher) Close() {
	// Producer won't write to our private delivery channel after Close/Flush
	close(p.publishResultChan)
}
