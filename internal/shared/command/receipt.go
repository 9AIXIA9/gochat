package command

import (
	"gochat/internal/shared/kernel"
	"time"
)

type ReceiptID string

func (id ReceiptID) String() string {
	return string(id)
}

type Receipt interface {
	ID() ReceiptID
	AggregateID() kernel.ID
	Action() Action
	Status() ReceiptStatus
	OccurredAt() time.Time
	Payload() []byte
	Headers() map[string]string
	AddHeader(key string, value string)
	AddHeaders(headers map[string]string)
}
