package websocket

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicUserSessionEnded event.Topic = "user_session.ended"

var _ event.SpecificEvent = (*UserSessionEndedEvent)(nil)

type UserSessionEndedEvent struct {
	*event.StandardEvent
}

// func ToUserSessionEndedEvent(ev event.Event) (*UserSessionEndedEvent, error) {
func _(ev event.Event) (*UserSessionEndedEvent, error) {
	e := &UserSessionEndedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUserSessionEndedEvent(
	userID kernel.UserID,
	generator event.IDGenerator,
) (*UserSessionEndedEvent, error) {
	e := &UserSessionEndedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicUserSessionEnded, payload, generator)
	return e, nil
}

func (e *UserSessionEndedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UserSessionEndedEvent) Unmarshal([]byte) error {
	return nil
}
