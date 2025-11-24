package websocket

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicUserSessionStarted event.Topic = "user_session.started"

var _ event.SpecificEvent = (*UserSessionStartedEvent)(nil)

type UserSessionStartedEvent struct {
	*event.StandardEvent
}

func ToUserSessionStartedEvent(ev event.Event) (*UserSessionStartedEvent, error) {
	e := &UserSessionStartedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUserSessionStartedEvent(
	userID kernel.UserID,
	generator event.IDGenerator,
) (*UserSessionStartedEvent, error) {
	e := &UserSessionStartedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicUserSessionStarted, payload, generator)
	return e, nil
}

func (e *UserSessionStartedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UserSessionStartedEvent) Unmarshal([]byte) error {
	return nil
}
