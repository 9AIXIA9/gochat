package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicUserCreated event.Topic = "notification.user.created"

var _ event.SpecificEvent = (*UserCreatedEvent)(nil)

type UserCreatedEvent struct {
	email  kernel.Email
	number UserNumber
	*event.StandardEvent
}

func ToUserCreatedEvent(ev event.Event) (*UserCreatedEvent, error) {
	e := &UserCreatedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUserCreatedEvent(id event.ID, userID kernel.UserID, email kernel.Email, number UserNumber) (*UserCreatedEvent, error) {
	e := &UserCreatedEvent{
		email:  email,
		number: number,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicUserCreated, payload)
	return e, nil
}

func (e *UserCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Email  kernel.Email
		Number UserNumber
	}
	return json.Marshal(&Alias{
		Email:  e.email,
		Number: e.number,
	})
}

func (e *UserCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Email  kernel.Email
		Number UserNumber
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.email = tmp.Email
	return nil
}

func (e *UserCreatedEvent) Number() UserNumber {
	return e.number
}

func (e *UserCreatedEvent) Email() kernel.Email {
	return e.email
}
