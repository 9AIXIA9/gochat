package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

var _ event.SpecificEvent = (*UserLoggedInEvent)(nil)

type UserLoggedInEvent struct {
	*event.StandardEvent
}

//func ToUserLoggedInEvent(ev event.Event) (*UserLoggedInEvent, error) {
//	e := &UserLoggedInEvent{StandardEvent: event.NewStandardEventFrom(ev)}
//	if len(ev.Payload()) > 0 {
//		if err := e.Unmarshal(ev.Payload()); err != nil {
//			return nil, err
//		}
//	}
//	return e, nil
//}

func NewUserLoggedInEvent(id event.ID, userID kernel.UserID) (*UserLoggedInEvent, error) {
	e := &UserLoggedInEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicUserLoggedIn, payload)
	return e, nil
}

func (e *UserLoggedInEvent) Marshal() ([]byte, error) {
	// no payload for this event; return an empty JSON object for consistency
	return []byte("{}"), nil
}

func (e *UserLoggedInEvent) Unmarshal(_ []byte) error {
	// nothing to unmarshal; accept empty or any JSON object
	return nil
}
