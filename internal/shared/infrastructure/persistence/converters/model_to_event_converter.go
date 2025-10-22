package converters

import (
	"gochat/internal/shared/errors"
	"gochat/internal/shared/infrastructure/persistence/models"
	kernelEvent "gochat/internal/shared/kernel/event"
	"time"
)

type ModelToSpecificEventConverter[SpecificEvent kernelEvent.Event] struct {
	factory func(id kernelEvent.ID, topic kernelEvent.Topic, occurredAt time.Time, payload []byte) (SpecificEvent, error)
}

func (c *ModelToSpecificEventConverter[SpecificEvent]) ToModels(gormEvents []*models.Event) ([]SpecificEvent, error) {
	if len(gormEvents) == 0 {
		return nil, errors.ErrEmptyInput
	}

	events := make([]SpecificEvent, 0, len(gormEvents))
	for _, gormEvent := range gormEvents {
		event, err := c.factory(gormEvent.ID, gormEvent.Topic, gormEvent.OccurredAt, gormEvent.Payload)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	return events, nil
}
