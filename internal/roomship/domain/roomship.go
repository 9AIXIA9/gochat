package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomshipRole string

const (
	OwnerRole  RoomshipRole = "owner"
	MemberRole RoomshipRole = "member"
)

type RoomshipID kernel.ID

func (r RoomshipID) String() string {
	return string(r)
}

type Roomship struct {
	id        RoomshipID
	roomID    kernel.RoomID
	userID    kernel.UserID
	role      RoomshipRole
	createdAt time.Time //UTC

	manager *event.Manager
}

func LoadRoomship(
	id RoomshipID,
	roomID kernel.RoomID,
	userID kernel.UserID,
	role RoomshipRole,
	createdAt time.Time,
) *Roomship {
	return &Roomship{
		id:        id,
		roomID:    roomID,
		userID:    userID,
		role:      role,
		createdAt: createdAt,
		manager:   event.NewEventManager(),
	}
}

func CreateRoomship(
	userID kernel.UserID,
	roomID kernel.RoomID,
	role RoomshipRole,
	roomshipIDGenerator RoomshipIDGenerator,
	idGenerator event.IDGenerator,
) *Roomship {
	roomship := &Roomship{
		id:        roomshipIDGenerator.Generate(),
		roomID:    roomID,
		userID:    userID,
		role:      role,
		createdAt: time.Now().UTC(),
		manager:   event.NewEventManager(),
	}

	ev, err := NewRoomshipCreatedEvent(roomship.id, idGenerator)
	if err != nil {
		return nil
	}

	roomship.manager.RecordEvent(ev)
	return roomship
}

func (r *Roomship) ID() RoomshipID {
	return r.id
}

func (r *Roomship) RoomID() kernel.RoomID {
	return r.roomID
}

func (r *Roomship) UserID() kernel.UserID {
	return r.userID
}

func (r *Roomship) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Roomship) Role() RoomshipRole {
	return r.role
}

func (r *Roomship) GetEvents() []event.Event {
	return r.manager.GetEvents()
}
