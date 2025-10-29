package event

import (
	"gochat/internal/shared/kernel"
	"time"
)

var _ Event = (*StandardEvent)(nil)

type StandardEvent struct {
	id          ID
	aggregateID kernel.ID
	occurredAt  time.Time //UTC
	topic       Topic
	payload     []byte
}

func NewStandardEvent(id ID, aggregateID kernel.ID, occurredAt time.Time, topic Topic, payload []byte) *StandardEvent {
	return &StandardEvent{
		id:          id,
		aggregateID: aggregateID,
		occurredAt:  occurredAt.UTC(),
		topic:       topic,
		payload:     payload,
	}
}

func (e *StandardEvent) ID() ID {
	return e.id
}

func (e *StandardEvent) AggregateID() kernel.ID {
	return e.aggregateID
}

func (e *StandardEvent) Topic() Topic {
	return e.topic
}

func (e *StandardEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *StandardEvent) Payload() []byte {
	return e.payload
}
