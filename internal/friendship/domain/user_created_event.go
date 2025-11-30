package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicUserCreated event.Topic = "friendship.user.created"

var _ event.SpecificEvent = (*UserCreatedEvent)(nil)

type UserCreatedEvent struct {
	number kernel.UserNumber
	*event.StandardEvent
}

func ToUserCreatedEvent(ev event.Event) (*UserCreatedEvent, error) {
	if ev.Topic() != TopicUserCreated {
		return nil, myErrors.ErrWrongEventType
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
	generator event.IDGenerator,
) (*UserCreatedEvent, error) {
	e := &UserCreatedEvent{
		number: number,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicUserCreated, payload, generator)
	return e, nil
}

func (e *UserCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Number kernel.UserNumber
	}
	return json.Marshal(&Alias{
		Number: e.number,
	})
}

func (e *UserCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Number kernel.UserNumber
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.number = tmp.Number
	return nil
}

func (e *UserCreatedEvent) Number() kernel.UserNumber {
	return e.number
}
