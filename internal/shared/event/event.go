package event

import (
	"gochat/internal/shared/kernel"
	"time"
)

type ID kernel.ID

type Topic string

type Event interface {
	ID() ID
	AggregateID() kernel.ID
	Topic() Topic
	OccurredAt() time.Time
	Payload() []byte
}

type SpecificEvent interface {
	Event
	kernel.Serializer
}
