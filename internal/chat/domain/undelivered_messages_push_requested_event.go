package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicUndeliveredMessagesPushRequested event.Topic = "chat.undelivered_messages_push.requested"

var _ event.SpecificEvent = (*UndeliveredMessagesPushRequestedEvent)(nil)

type UndeliveredMessagesPushRequestedEvent struct {
	*event.StandardEvent
}

func ToUndeliveredMessagesPushRequestedEvent(ev event.Event) (*UndeliveredMessagesPushRequestedEvent, error) {
	if ev.Topic() != TopicUndeliveredMessagesPushRequested {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &UndeliveredMessagesPushRequestedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUndeliveredMessagesPushRequestedEvent(
	userID kernel.UserID,
	generator event.IDGenerator,
) (*UndeliveredMessagesPushRequestedEvent, error) {
	e := &UndeliveredMessagesPushRequestedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicUndeliveredMessagesPushRequested, payload, generator)
	return e, nil
}

func (e *UndeliveredMessagesPushRequestedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UndeliveredMessagesPushRequestedEvent) Unmarshal([]byte) error {
	return nil
}
