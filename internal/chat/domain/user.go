package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type User struct {
	id              kernel.UserID
	privateMessages []*PrivateMessage
	roomMessages    []*RoomMessage

	eventManager *event.Manager
}

func NewUser(id kernel.UserID, privateMessages []*PrivateMessage, roomMessages []*RoomMessage) *User {
	return &User{
		id:              id,
		privateMessages: privateMessages,
		roomMessages:    roomMessages,
		eventManager:    event.NewEventManager(),
	}
}

func (u *User) SendPrivateMessage(message *PrivateMessage, generator event.IDGenerator) error {
	u.privateMessages = append(u.privateMessages, message)
	ev, err := NewPrivateMessageCreatedEvent(
		generator.Generate(),
		message.id,
		message.sender,
		message.recipient,
		message.content,
		message.sentAt,
	)
	if err != nil {
		return err
	}
	u.eventManager.RecordEvent(ev)
	return nil
}

func (u *User) SendRoomMessage(message *RoomMessage, generator event.IDGenerator) error {
	u.roomMessages = append(u.roomMessages, message)

	ev, err := NewRoomMessageCreatedEvent(
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

func (u *User) PrivateMessages() []*PrivateMessage {
	return u.privateMessages
}
func (u *User) RoomMessages() []*RoomMessage {
	return u.roomMessages
}

func (u *User) GetEvents() []event.Event {
	return u.eventManager.GetEvents()
}
