//go:generate mockgen -source=command.go -destination=./mocks/mock_command.go -package=mocks
package command

import (
	"gochat/internal/shared/kernel"
	"time"
)

type ID string

func (i ID) String() string {
	return string(i)
}

type Action string

func (t Action) String() string {
	return string(t)
}

type Command interface {
	ID() ID
	AggregateID() kernel.ID
	Action() Action
	OccurredAt() time.Time
	Payload() []byte
	Headers() map[string]string
	AddHeader(key string, value string)
	AddHeaders(headers map[string]string)
}

type SpecificCommand interface {
	Command
	kernel.Serializer
}
