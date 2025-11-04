package event

import (
	"gochat/internal/shared/kernel"
	"time"
)

type ID kernel.ID

func (i ID) String() string {
	return kernel.ID(i).String()
}

type Topic string

func (t Topic) String() string {
	return string(t)
}

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
