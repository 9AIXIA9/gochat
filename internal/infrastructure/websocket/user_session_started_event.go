package websocket

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicUserSessionStarted event.Topic = "user_session.started"

var _ event.SpecificEvent = (*UserSessionStartedEvent)(nil)

type UserSessionStartedEvent struct {
	*event.StandardEvent
}

func ToUserSessionStartedEvent(ev event.Event) (*UserSessionStartedEvent, error) {
	e := &UserSessionStartedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUserSessionStartedEvent(
	id event.ID,
	userID kernel.UserID,
) (*UserSessionStartedEvent, error) {
	e := &UserSessionStartedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicUserSessionStarted, payload)
	return e, nil
}

func (e *UserSessionStartedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UserSessionStartedEvent) Unmarshal([]byte) error {
	return nil
}
