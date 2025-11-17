package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type User struct {
	id     kernel.UserID
	number kernel.UserNumber

	eventManager *event.Manager
}

func NewUser(id kernel.UserID, number kernel.UserNumber) *User {
	return &User{
		id:           id,
		number:       number,
		eventManager: event.NewEventManager(),
	}
}

func (u *User) ReceiveMessage(id kernel.MessageID, sender kernel.UserID, content string, generator event.IDGenerator) (*Message, error) {
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

func (u *User) Number() kernel.UserNumber {
	return u.number
}

func (u *User) GetEvents() []event.Event {
	return u.eventManager.GetEvents()
}
