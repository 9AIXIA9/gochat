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

	members []kernel.UserID

	eventManager *event.Manager
}

func NewRoom(id RoomID, number RoomNumber, members []kernel.UserID) *Room {
	return &Room{
		id:           id,
		number:       number,
		members:      members,
		eventManager: event.NewEventManager(),
	}
}

func (r *Room) ReceiveMessage(id MessageID, sender kernel.UserID, content string, generator event.IDGenerator) (*Message, error) {
	for i, member := range r.members {
		if member == sender {
			now := time.Now().UTC()

			ev, err := NewRoomMessageCreatedEvent(
				generator.Generate(),
				id,
				r.id,
				append(r.members[:i], r.members[i+1:]...),
				sender,
				content,
				now,
			)
			if err != nil {
				return nil, err
			}
			r.eventManager.RecordEvent(ev)
			return NewMessage(id, sender, content, now), nil
		}
	}
	return nil, myErrors.ErrNotBelongTo
}

func (r *Room) ID() RoomID {
	return r.id
}

func (r *Room) Number() RoomNumber {
	return r.number
}

func (r *Room) Members() []kernel.UserID {
	return r.members
}

func (r *Room) GetEvents() []event.Event {
	return r.eventManager.GetEvents()
}
