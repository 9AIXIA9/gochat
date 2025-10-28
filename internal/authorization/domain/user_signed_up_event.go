package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

var _ event.SpecificEvent = (*UserSignedUpEvent)(nil)

type UserSignedUpEvent struct {
	email kernel.Email
	*event.StandardEvent
}

//func ToUserSignedUpEvent(ev event.Event) (*UserSignedUpEvent, error) {
//	e := &UserSignedUpEvent{StandardEvent: event.NewStandardEventFrom(ev)}
//	if len(ev.Payload()) > 0 {
//		if err := e.Unmarshal(ev.Payload()); err != nil {
//			return nil, err
//		}
//	}
//	return e, nil
//}

func NewUserSignedUpEvent(id event.ID, userID kernel.UserID, email kernel.Email) (*UserSignedUpEvent, error) {
	e := &UserSignedUpEvent{email: email}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicUserSignedUp, payload)
	return e, nil
}

func (e *UserSignedUpEvent) Email() kernel.Email {
	return e.email
}

func (e *UserSignedUpEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Email kernel.Email
	}
	return json.Marshal(Alias{Email: e.email})
}

func (e *UserSignedUpEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Email kernel.Email
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.email = tmp.Email
	return nil
}
