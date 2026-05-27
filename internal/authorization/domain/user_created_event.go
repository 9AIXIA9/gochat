package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicUserCreated event.Topic = "user.created"

var _ event.SpecificEvent = (*UserCreatedEvent)(nil)

type UserCreatedEvent struct {
	number   kernel.UserNumber
	email    kernel.Email
	signedAt time.Time
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

func NewUserCreatedEvent(
	userID kernel.UserID,
	number kernel.UserNumber,
	email kernel.Email,
	signedAt time.Time,
	generator event.IDGenerator,
) (*UserCreatedEvent, error) {
	e := &UserCreatedEvent{
		number:   number,
		email:    email,
		signedAt: signedAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicUserCreated, payload, generator)
	return e, nil
}

func (e *UserCreatedEvent) Number() kernel.UserNumber {
	return e.number
}

func (e *UserCreatedEvent) Email() kernel.Email {
	return e.email
}

func (e *UserCreatedEvent) SignedAt() time.Time {
	return e.signedAt
}

func (e *UserCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Number   kernel.UserNumber
		Email    kernel.Email
		SignedAt time.Time
	}
	return json.Marshal(Alias{
		Number:   e.number,
		Email:    e.email,
		SignedAt: e.signedAt,
	})
}

func (e *UserCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Number   kernel.UserNumber
		Email    kernel.Email
		SignedAt time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.signedAt = tmp.SignedAt
	e.email = tmp.Email
	e.number = tmp.Number
	return nil
}
