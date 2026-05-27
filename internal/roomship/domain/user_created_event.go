package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
)

const TopicUserCreated event.Topic = "user.created"

var _ event.SpecificEvent = (*UserCreatedEvent)(nil)

type UserCreatedEvent struct {
	*event.StandardEvent
}

func ToUserCreatedEvent(ev event.Event) (*UserCreatedEvent, error) {
	if ev.Topic() != TopicUserCreated {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &UserCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *UserCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UserCreatedEvent) Unmarshal(_ []byte) error {
	return nil
}
