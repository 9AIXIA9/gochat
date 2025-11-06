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
	id      RoomID
	number  RoomNumber
	members []*User
}

func NewRoom(id RoomID, number RoomNumber, members []*User) *Room {
	return &Room{
		id:      id,
		number:  number,
		members: members,
	}
}

func (r *Room) ReceiveMessage(id MessageID, sender kernel.UserID, content string, generator event.IDGenerator) error {
	if !r.IsMember(sender) {
		return myErrors.ErrNotBelongTo
	}

	for _, member := range r.members {
		if member.id != sender {
			if err := member.ReceiveMessage(id, sender, content, time.Now().UTC(), generator); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Room) IsMember(id kernel.UserID) bool {
	for _, member := range r.members {
		if member.id == id {
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

func (r *Room) Members() []*User {
	return r.members
}

func (r *Room) GetEvents() []event.Event {
	var events []event.Event
	for _, member := range r.members {
		events = append(events, member.GetEvents()...)
	}
	return events
}
