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
	headers     map[string]string
}

func LoadStandardEvent(
	id ID,
	aggregateID kernel.ID,
	occurredAt time.Time,
	topic Topic,
	payload []byte,
	headers map[string]string,
) *StandardEvent {
	return &StandardEvent{
		id:          id,
		aggregateID: aggregateID,
		occurredAt:  occurredAt,
		topic:       topic,
		payload:     payload,
		headers:     headers,
	}
}

func LoadStandardEventFromEvent(e Event) *StandardEvent {
	if se, ok := e.(*StandardEvent); ok {
		return se
	}
	return &StandardEvent{
		id:          e.ID(),
		aggregateID: e.AggregateID(),
		occurredAt:  e.OccurredAt(),
		topic:       e.Topic(),
		payload:     e.Payload(),
		headers:     e.Headers(),
	}
}

func NewStandardEvent(
	aggregateID kernel.ID,
	topic Topic,
	payload []byte,
	generator IDGenerator,
) *StandardEvent {
	return &StandardEvent{
		id:          generator.Generate(),
		aggregateID: aggregateID,
		occurredAt:  time.Now().UTC(),
		topic:       topic,
		payload:     payload,
	}
}

func (e *StandardEvent) AddHeader(key string, value string) {
	if e.headers == nil {
		e.headers = make(map[string]string)
	}
	e.headers[key] = value
}

func (e *StandardEvent) AddHeaders(headers map[string]string) {
	if e.headers == nil {
		e.headers = make(map[string]string)
	}
	for k, v := range headers {
		e.headers[k] = v
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

func (e *StandardEvent) Headers() map[string]string {
	return e.headers
}
