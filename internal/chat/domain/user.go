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

	messagesReceived []*Message

	eventManager *event.Manager
}

func NewUser(id kernel.UserID, number UserNumber, messagesReceived []*Message) *User {
	return &User{
		id:               id,
		number:           number,
		messagesReceived: messagesReceived,
		eventManager:     event.NewEventManager(),
	}
}

func (u *User) ReceiveMessage(id MessageID, sender kernel.UserID, content string, sentAt time.Time, generator event.IDGenerator) error {
	state := NewMessage(id, MessageStateReceived, sender, content, sentAt)

	u.messagesReceived = append(u.messagesReceived, state)

	ev, err := NewPrivateMessageReceivedEvent(generator.Generate(), id, u.id, sender, content, sentAt)
	if err != nil {
		return err
	}

	u.eventManager.RecordEvent(ev)
	return nil
}

func (u *User) ID() kernel.UserID {
	return u.id
}

func (u *User) Number() UserNumber {
	return u.number
}

func (u *User) MessagesReceived() []*Message {
	return u.messagesReceived
}

func (u *User) GetEvents() []event.Event {
	return u.eventManager.GetEvents()
}
