package converter

import (
	kernelEvent "gochat/internal/shared/event"
)

type EventInterfaceConverter struct {
}

func (c *EventInterfaceConverter) ToStandard(event kernelEvent.Event) *kernelEvent.StandardEvent {
	return kernelEvent.NewStandardEvent(event.ID(), event.AggregateID(), event.OccurredAt(), event.Topic(), event.Payload())
}

func (c *EventInterfaceConverter) ToStandards(events []kernelEvent.Event) []*kernelEvent.StandardEvent {
	domainEvents := make([]*kernelEvent.StandardEvent, len(events))
	for i, event := range events {
		domainEvents[i] = c.ToStandard(event)
	}
	return domainEvents
}

func (c *EventInterfaceConverter) ToInterface(event *kernelEvent.StandardEvent) kernelEvent.Event {
	return event
}
func (c *EventInterfaceConverter) ToInterfaces(events []*kernelEvent.StandardEvent) []kernelEvent.Event {
	modelEvents := make([]kernelEvent.Event, len(events))
	for i, event := range events {
		modelEvents[i] = c.ToInterface(event)
	}
	return modelEvents
}
