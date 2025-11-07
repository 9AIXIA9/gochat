package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomID kernel.ID

func (r RoomID) String() string {
	return string(r)
}

type RoomNumber kernel.Number

func (n RoomNumber) String() string {
	return string(n)
}

func (n RoomNumber) Validate() error {
	if len(n) == 0 {
		return myErrors.ErrInvalidNumber
	}
	return nil
}

type Room struct {
	id     RoomID
	number RoomNumber

	members          []kernel.UserID
	messagesReceived []*Message

	eventManager *event.Manager
}

func NewRoom(id RoomID, number RoomNumber, members []kernel.UserID) *Room {
	return &Room{
		id:               id,
		number:           number,
		members:          members,
		messagesReceived: make([]*Message, 0),
		eventManager:     event.NewEventManager(),
	}
}

func (r *Room) ReceiveMessage(id MessageID, sender kernel.UserID, content string, generator event.IDGenerator) error {
	if !r.IsMember(sender) {
		return myErrors.ErrNotBelongTo
	}

	r.messagesReceived = append(r.messagesReceived, NewMessage(id, MessageStateReceived, sender, content, time.Now().UTC()))

	ev, err := NewRoomMessageCreatedEvent(generator.Generate(), id, r.id)
	if err != nil {
		return err
	}
	r.eventManager.RecordEvent(ev)
	return nil
}

func (r *Room) IsMember(id kernel.UserID) bool {
	for _, member := range r.members {
		if member == id {
			return true
		}
	}
	return false
}

func (r *Room) ID() RoomID {
	return r.id
}

func (r *Room) Number() RoomNumber {
	return r.number
}

func (r *Room) MessagesReceived() []*Message {
	return r.messagesReceived
}

func (r *Room) Members() []kernel.UserID {
	return r.members
}

func (r *Room) GetEvents() []event.Event {
	return r.eventManager.GetEvents()
}
