package command

import (
	"encoding/json"
	"gochat/internal/shared/kernel"
	"time"
)

var _ Receipt = (*StandardReceipt)(nil)

type StandardReceipt struct {
	id          ReceiptID
	commandID   ID
	aggregateID kernel.ID
	occurredAt  time.Time //UTC
	action      Action
	status      ReceiptStatus
	headers     map[string]string
}

func LoadStandardReceipt(
	id ReceiptID,
	commandID ID,
	aggregateID kernel.ID,
	occurredAt time.Time,
	action Action,
	status ReceiptStatus,
	headers map[string]string,
) *StandardReceipt {
	return &StandardReceipt{
		id:          id,
		commandID:   commandID,
		aggregateID: aggregateID,
		occurredAt:  occurredAt,
		action:      action,
		status:      status,
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
		commandID:   command.ID(),
		aggregateID: command.AggregateID(),
		occurredAt:  time.Now().UTC(),
		action:      command.Action(),
		status:      status,
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

func (r *StandardReceipt) CommandID() ID {
	return r.commandID
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

func (r *StandardReceipt) Headers() map[string]string {
	return r.headers
}

func (r *StandardReceipt) Marshal() ([]byte, error) {
	type Alias struct {
		ID          ReceiptID
		CommandID   ID
		AggregateID kernel.ID
		OccurredAt  time.Time //UTC
		Action      Action
		Status      ReceiptStatus
	}
	return json.Marshal(&Alias{
		ID:          r.id,
		CommandID:   r.commandID,
		AggregateID: r.aggregateID,
		OccurredAt:  r.occurredAt,
		Action:      r.action,
		Status:      r.status,
	})
}

func (r *StandardReceipt) Unmarshal(data []byte) error {
	type Alias struct {
		ID          ReceiptID
		CommandID   ID
		AggregateID kernel.ID
		OccurredAt  time.Time //UTC
		Action      Action
		Status      ReceiptStatus
	}

	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	r.id = tmp.ID
	r.commandID = tmp.CommandID
	r.aggregateID = tmp.AggregateID
	r.occurredAt = tmp.OccurredAt
	r.action = tmp.Action
	r.status = tmp.Status
	return nil
}
