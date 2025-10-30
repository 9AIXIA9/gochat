package uuid

import (
	"gochat/internal/shared/event"

	"github.com/google/uuid"
)

var _ event.IDGenerator = (*EventIDGenerator)(nil)

type EventIDGenerator struct{}

func NewEventIDGenerator() *EventIDGenerator {
	return &EventIDGenerator{}
}

func (EventIDGenerator) Generate() event.ID {
	return event.ID(uuid.Must(uuid.NewV7()).String())
}
