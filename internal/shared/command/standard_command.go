package command

import (
	"gochat/internal/shared/kernel"
	"time"
)

var _ Command = (*StandardCommand)(nil)

type StandardCommand struct {
	id          ID
	aggregateID kernel.ID
	occurredAt  time.Time //UTC
	action      Action
	payload     []byte
	headers     map[string]string
}

func LoadStandardCommand(
	id ID,
	aggregateID kernel.ID,
	occurredAt time.Time,
	action Action,
	payload []byte,
	headers map[string]string,
) *StandardCommand {
	return &StandardCommand{
		id:          id,
		aggregateID: aggregateID,
		occurredAt:  occurredAt,
		action:      action,
		payload:     payload,
		headers:     headers,
	}
}

func LoadStandardCommandFromCommand(e Command) *StandardCommand {
	if se, ok := e.(*StandardCommand); ok {
		return se
	}
	return &StandardCommand{
		id:          e.ID(),
		aggregateID: e.AggregateID(),
		occurredAt:  e.OccurredAt(),
		action:      e.Action(),
		payload:     e.Payload(),
		headers:     e.Headers(),
	}
}

func NewStandardCommand(
	aggregateID kernel.ID,
	action Action,
	payload []byte,
	generator IDGenerator,
) *StandardCommand {
	return &StandardCommand{
		id:          generator.Generate(),
		aggregateID: aggregateID,
		occurredAt:  time.Now().UTC(),
		action:      action,
		payload:     payload,
	}
}

func (e *StandardCommand) AddHeader(key string, value string) {
	if e.headers == nil {
		e.headers = make(map[string]string)
	}
	e.headers[key] = value
}

func (e *StandardCommand) AddHeaders(headers map[string]string) {
	if e.headers == nil {
		e.headers = make(map[string]string)
	}
	for k, v := range headers {
		e.headers[k] = v
	}
}

func (e *StandardCommand) ID() ID {
	return e.id
}

func (e *StandardCommand) AggregateID() kernel.ID {
	return e.aggregateID
}

func (e *StandardCommand) Action() Action {
	return e.action
}

func (e *StandardCommand) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *StandardCommand) Payload() []byte {
	return e.payload
}

func (e *StandardCommand) Headers() map[string]string {
	return e.headers
}
