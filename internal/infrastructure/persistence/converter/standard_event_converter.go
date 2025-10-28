package converter

import (
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/model"
	kernelEvent "gochat/internal/shared/event"
)

var _ gormutils.GenericModelConverter[*model.Event, *kernelEvent.StandardEvent] = (*StandardEventConverter)(nil)

type StandardEventConverter struct {
}

func (c *StandardEventConverter) ToDomain(event *model.Event) *kernelEvent.StandardEvent {
	return kernelEvent.NewStandardEvent(event.ID, event.AggregateID, event.OccurredAt, event.Topic, event.Payload)
}

func (c *StandardEventConverter) ToModel(event *kernelEvent.StandardEvent) *model.Event {
	return &model.Event{
		ID:          event.ID(),
		AggregateID: event.AggregateID(),
		Topic:       event.Topic(),
		OccurredAt:  event.OccurredAt(),
		Payload:     event.Payload(),
	}
}
