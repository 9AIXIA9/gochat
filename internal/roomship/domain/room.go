package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomNumber kernel.Number

func (n RoomNumber) String() string {
	return string(n)
}

func (n RoomNumber) Validate() error {
	return kernel.Number(n).Validate()
}

type Room struct {
	id                kernel.RoomID
	number            RoomNumber
	ownerID           kernel.UserID
	passwordEncrypted PasswordEncrypted
	maxMemberCount    int
	createdAt         time.Time

	manager *event.Manager
}

func LoadRoom(
	id kernel.RoomID,
	ownerID kernel.UserID,
	number RoomNumber,
	passwordEncrypted PasswordEncrypted,
	maxMemberCount int,
	createdAt time.Time,
) *Room {
	return &Room{
		id:                id,
		number:            number,
		ownerID:           ownerID,
		passwordEncrypted: passwordEncrypted,
		maxMemberCount:    maxMemberCount,
		createdAt:         createdAt,
		manager:           event.NewEventManager(),
	}
}

func CreateRoom(
	ownerID kernel.UserID,
	maxMemberCount int,
	passwordEncrypted PasswordEncrypted,
	roomIDGenerator RoomIDGenerator,
	roomNumberGenerator RoomNumberGenerator,
	idGenerator event.IDGenerator,
) (*Room, error) {
	if maxMemberCount < 2 {
		return nil, ErrInvalidMaxMemberCount
	}

	room := &Room{
		id:                roomIDGenerator.Generate(),
		number:            roomNumberGenerator.Generate(),
		ownerID:           ownerID,
		passwordEncrypted: passwordEncrypted,
		maxMemberCount:    maxMemberCount,
		createdAt:         time.Now().UTC(),
		manager:           event.NewEventManager(),
	}

	ev, err := NewRoomCreatedEvent(room.id, idGenerator)
	if err != nil {
		return nil, err
	}

	room.manager.RecordEvent(ev)

	return room, nil
}

func (r *Room) ID() kernel.RoomID {
	return r.id
}

func (r *Room) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Room) MaxMemberCount() int {
	return r.maxMemberCount
}

func (r *Room) OwnerID() kernel.UserID {
	return r.ownerID
}

func (r *Room) Number() RoomNumber {
	return r.number
}

func (r *Room) PasswordEncrypted() PasswordEncrypted {
	return r.passwordEncrypted
}

func (r *Room) GetEvents() []event.Event {
	return r.manager.GetEvents()
}
