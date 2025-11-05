package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
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

func (u *User) SendMessage(id MessageID, recipient kernel.UserID, content string, generator event.IDGenerator) error {
	state := NewRecipientMessageState(recipient, MessageStateCreated)
	states := []*RecipientMessageState{state}
	message := NewMessage(id, u.id, content, time.Now().UTC(), states)

	u.messages = append(u.messages, message)
	ev, err := NewMessageCreatedEvent(
		generator.Generate(),
		message.id,
		message.sender,
		message.Recipients(),
		message.content,
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

func (u *User) Messages() []*Message {
	return u.messages
}

func (u *User) GetEvents() []event.Event {
	return u.eventManager.GetEvents()
}
