package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type UserNumber kernel.Number

func (n UserNumber) String() string {
	return string(n)
}

func (n UserNumber) Validate() error {
	if len(n) == 0 {
		return myErrors.ErrInvalidNumber
	}
	return nil
}

type User struct {
	id     kernel.UserID
	number UserNumber

	eventManager *event.Manager
}

func NewUser(id kernel.UserID, number UserNumber) *User {
	return &User{
		id:           id,
		number:       number,
		eventManager: event.NewEventManager(),
	}
}

func (u *User) ReceiveMessage(id MessageID, sender kernel.UserID, content string, generator event.IDGenerator) (*Message, error) {
	now := time.Now().UTC()
	message := NewMessage(id, sender, content, now)

	ev, err := NewPrivateMessageCreatedEvent(generator.Generate(), id, u.id, sender, content, now)
	if err != nil {
		return nil, err
	}

	u.eventManager.RecordEvent(ev)
	return message, nil
}

func (u *User) ID() kernel.UserID {
	return u.id
}

func (u *User) Number() UserNumber {
	return u.number
}

func (u *User) GetEvents() []event.Event {
	return u.eventManager.GetEvents()
}
