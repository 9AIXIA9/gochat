package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicUserConnected event.Topic = "notification.user.connected"

var _ event.SpecificEvent = (*UserConnectedEvent)(nil)

type UserConnectedEvent struct {
	*event.StandardEvent
}

func ToUserConnectedEvent(ev event.Event) (*UserConnectedEvent, error) {
	e := &UserConnectedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUserConnectedEvent(id event.ID, userID kernel.UserID) (*UserConnectedEvent, error) {
	e := &UserConnectedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicUserConnected, payload)
	return e, nil
}

func (e *UserConnectedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UserConnectedEvent) Unmarshal([]byte) error {
	return nil
}
