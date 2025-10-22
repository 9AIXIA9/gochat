package converters

import (
	"gochat/internal/shared/errors"
	"gochat/internal/shared/infrastructure/persistence/models"
	kernelEvent "gochat/internal/shared/kernel/event"
)

type EventToModelConverter struct {
}

func (c *EventToModelConverter) ToModels(events []kernelEvent.Event) ([]*models.Event, error) {
	if len(events) == 0 {
		return nil, errors.ErrEmptyInput
	}

	gormEvents := make([]*models.Event, 0, len(events))
	for _, event := range events {
		payload, err := event.Marshal()
		if err != nil {
			return nil, err
		}
		gormEvents = append(gormEvents, &models.Event{
			ID:         event.ID(),
			Topic:      event.Topic(),
			OccurredAt: event.OccurredAt(),
			Payload:    payload,
		})
	}
	return gormEvents, nil
}
