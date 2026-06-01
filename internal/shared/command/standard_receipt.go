package command

import (
	"gochat/internal/shared/kernel"
	"time"
)

var _ Receipt = (*StandardReceipt)(nil)

type StandardReceipt struct {
	id          ReceiptID
	aggregateID kernel.ID
	occurredAt  time.Time //UTC
	action      Action
	status      ReceiptStatus
	payload     []byte
	headers     map[string]string
}

func LoadStandardReceipt(
	id ReceiptID,
	aggregateID kernel.ID,
	occurredAt time.Time,
	action Action,
	status ReceiptStatus,
	payload []byte,
	headers map[string]string,
) *StandardReceipt {
	return &StandardReceipt{
		id:          id,
		aggregateID: aggregateID,
		occurredAt:  occurredAt,
		action:      action,
		status:      status,
		payload:     payload,
		headers:     headers,
	}
}

func NewStandardReceiptFromCommand(
	command Command,
	status ReceiptStatus,
	generator ReceiptIDGenerator,
) *StandardReceipt {
	return &StandardReceipt{
		id:          generator.Generate(),
		aggregateID: command.AggregateID(),
		occurredAt:  time.Now().UTC(),
		action:      command.Action(),
		status:      status,
		payload:     nil,
		headers:     command.Headers(),
	}
}

func (r *StandardReceipt) AddHeader(key string, value string) {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[key] = value
}

func (r *StandardReceipt) AddHeaders(headers map[string]string) {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	for k, v := range headers {
		r.headers[k] = v
	}
}

func (r *StandardReceipt) ID() ReceiptID {
	return r.id
}

func (r *StandardReceipt) AggregateID() kernel.ID {
	return r.aggregateID
}

func (r *StandardReceipt) Action() Action {
	return r.action
}

func (r *StandardReceipt) Status() ReceiptStatus {
	return r.status
}

func (r *StandardReceipt) OccurredAt() time.Time {
	return r.occurredAt
}

func (r *StandardReceipt) Payload() []byte {
	return r.payload
}

func (r *StandardReceipt) Headers() map[string]string {
	return r.headers
}
