package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicPrivateMessageCreated event.Topic = "chat.private_message.created"

var _ event.SpecificEvent = (*PrivateMessageCreatedEvent)(nil)

type PrivateMessageCreatedEvent struct {
	*event.StandardEvent
}

func ToPrivateMessageCreatedEvent(ev event.Event) (*PrivateMessageCreatedEvent, error) {
	if ev.Topic() != TopicPrivateMessageCreated {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &PrivateMessageCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewPrivateMessageCreatedEvent(
	messageID kernel.MessageID,
	generator event.IDGenerator,
) (*PrivateMessageCreatedEvent, error) {
	e := &PrivateMessageCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(messageID), TopicPrivateMessageCreated, payload, generator)
	return e, nil
}

func (e *PrivateMessageCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *PrivateMessageCreatedEvent) Unmarshal([]byte) error {
	return nil
}
