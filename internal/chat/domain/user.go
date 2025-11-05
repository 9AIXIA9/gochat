package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type User struct {
	id              kernel.UserID
	privateMessages []*PrivateMessage

	eventManager *event.Manager
}

func NewUser(id kernel.UserID, privateMessages []*PrivateMessage) *User {
	return &User{
		id:              id,
		privateMessages: privateMessages,
		eventManager:    event.NewEventManager(),
	}
}

func (u *User) SendPrivateMessage(message *PrivateMessage, generator event.IDGenerator) error {
	u.privateMessages = append(u.privateMessages, message)
	ev, err := NewMessageCreatedEvent(
		generator.Generate(),
		message.id,
		message.sender,
		kernel.ID(message.recipient),
		message.content,
		message.mType,
		message.sentAt,
	)
	if err != nil {
		return err
	}
	u.eventManager.RecordEvent(ev)
	return nil
}

func (u *User) ID() kernel.UserID {
	return u.id
}

func (u *User) PrivateMessages() []*PrivateMessage {
	return u.privateMessages
}

func (u *User) GetEvents() []event.Event {
	return u.eventManager.GetEvents()
}
