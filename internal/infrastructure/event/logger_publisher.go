package event

import (
	"context"
	kEvent "gochat/internal/shared/event"

	"go.uber.org/zap"
)

// LoggerPublisher is a simple Publisher that logs events. Replace with Kafka/RabbitMQ etc in production.
type LoggerPublisher struct{}

func NewLoggerPublisher() *LoggerPublisher { return &LoggerPublisher{} }

var _ kEvent.Publisher = (*LoggerPublisher)(nil)

func (p *LoggerPublisher) PublishEvents(_ context.Context, events []kEvent.Event) error {
	for _, e := range events {
		zap.L().Info("publish event", zap.String("topic", string(e.Topic())), zap.String("id", string(e.ID())), zap.ByteString("payload", e.Payload()))
	}
	return nil
}
