package websocket

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicUserSessionEnded event.Topic = "user_session.ended"

var _ event.SpecificEvent = (*UserSessionEndedEvent)(nil)

type UserSessionEndedEvent struct {
	*event.StandardEvent
}

// func ToUserSessionEndedEvent(ev event.Event) (*UserSessionEndedEvent, error) {
func _(ev event.Event) (*UserSessionEndedEvent, error) {
	e := &UserSessionEndedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUserSessionEndedEvent(
	id event.ID,
	userID kernel.UserID,
) (*UserSessionEndedEvent, error) {
	e := &UserSessionEndedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicUserSessionEnded, payload)
	return e, nil
}

func (e *UserSessionEndedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UserSessionEndedEvent) Unmarshal([]byte) error {
	return nil
}
