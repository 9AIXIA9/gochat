package event

import (
	"gochat/internal/shared/kernel"
	"time"
)

var _ StandardEvent = (*standardEvent)(nil)

type ID kernel.ID

type StandardEvent interface {
	ID() ID
	Topic() Topic
	OccurredAt() time.Time
}

type standardEvent struct {
	id         ID
	occurredAt time.Time //UTC
	topic      Topic
}

func NewStandardEvent(id ID, occurredAt time.Time, topic Topic) StandardEvent {
	return &standardEvent{
		id:         id,
		occurredAt: occurredAt,
		topic:      topic,
	}
}

func CreateStandardEvent(topic Topic, generator kernel.IDGenerator) StandardEvent {
	return &standardEvent{
		id:         ID(generator.Generate()),
		occurredAt: time.Now().UTC(),
		topic:      topic,
	}
}

func (e *standardEvent) ID() ID {
	return e.id
}

func (e *standardEvent) Topic() Topic {
	return e.topic
}

func (e *standardEvent) OccurredAt() time.Time {
	return e.occurredAt
}
