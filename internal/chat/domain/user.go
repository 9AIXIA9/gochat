package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type User struct {
	id       kernel.UserID
	messages []*Message

	eventManager *event.Manager
}

func NewUser(id kernel.UserID, messages []*Message) *User {
	return &User{
		id:           id,
		messages:     messages,
		eventManager: event.NewEventManager(),
	}
}

func (u *User) SendMessage(message *Message, generator event.IDGenerator) error {
	u.messages = append(u.messages, message)
	ev, err := NewMessageSentEvent(
		generator.Generate(),
		message.sender,
		message.recipient,
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

func (u *User) ReceiveMessage(message *Message) {
	u.messages = append(u.messages, message)
}

func (u *User) ID() kernel.UserID {
	return u.id
}

func (u *User) Messages() []*Message {
	return u.messages
}

func (u *User) GetEvents() []event.Event {
	return u.eventManager.GetEvents()
}
