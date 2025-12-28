//go:generate mockgen -source=event.go -destination=./mocks/mock_event.go -package=mocks
package event

import (
	"gochat/internal/shared/kernel"
	"time"
)

type ID string

func (i ID) String() string {
	return string(i)
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
	Headers() map[string]string
	AddHeader(key string, value string)
	AddHeaders(headers map[string]string)
}

type SpecificEvent interface {
	Event
	kernel.Serializer
}
