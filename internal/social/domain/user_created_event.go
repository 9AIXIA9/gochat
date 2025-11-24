package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicUserCreated event.Topic = "social.user.created"

var _ event.SpecificEvent = (*UserCreatedEvent)(nil)

type UserCreatedEvent struct {
	*event.StandardEvent
}

func ToUserCreatedEvent(ev event.Event) (*UserCreatedEvent, error) {
	e := &UserCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUserCreatedEvent(userID kernel.UserID, generator event.IDGenerator) (*UserCreatedEvent, error) {
	e := &UserCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicUserCreated, payload, generator)
	return e, nil
}

func (e *UserCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UserCreatedEvent) Unmarshal(_ []byte) error {
	return nil
}
