package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomID kernel.ID

func (i RoomID) String() string {
	return string(i)
}

type RoomNumber kernel.Number

func (n RoomNumber) String() string {
	return string(n)
}

type Room struct {
	id                RoomID
	owner             kernel.UserID
	number            RoomNumber
	passwordEncrypted string
	members           []kernel.UserID
	memberCount       int
	maxMemberCount    int
	createdAt         time.Time
	eventManager      *event.Manager
}

func NewRoom(
	id RoomID,
	owner kernel.UserID,
	number RoomNumber,
	passwordEncrypted string,
	members []kernel.UserID,
	memberCount int,
	maxMemberCount int,
	createdAt time.Time,
) *Room {
	return &Room{
		id:                id,
		owner:             owner,
		number:            number,
		passwordEncrypted: passwordEncrypted,
		members:           members,
		memberCount:       memberCount,
		maxMemberCount:    maxMemberCount,
		createdAt:         createdAt,
		eventManager:      event.NewEventManager(),
	}
}

func (r *Room) Create(generator event.IDGenerator) error {
	ev, err := NewRoomCreatedEvent(generator.Generate(), r.id)
	if err != nil {
		return err
	}
	r.eventManager.RecordEvent(ev)
	return nil
}

func (r *Room) Join(userID kernel.UserID, generator event.IDGenerator) error {
	if r.maxMemberCount <= r.memberCount {
		return myErrors.ErrExceedMaxValue
	}

	if r.IsMember(userID) {
		return myErrors.ErrHasBeenDone
	}

	r.members = append(r.members, userID)
	r.memberCount++

	ev, err := NewRoomJoinedEvent(generator.Generate(), userID, r.id)
	if err != nil {
		return err
	}

	r.eventManager.RecordEvent(ev)
	return nil
}

func (r *Room) Leave(userID kernel.UserID, generator event.IDGenerator) error {
	if userID == r.owner {
		return myErrors.ErrOwnerCantLeave
	}

	for i, member := range r.members {
		if member == userID {
			r.members = append(r.members[:i], r.members[i+1:]...)
			r.memberCount--

			ev, err := NewRoomLeftEvent(generator.Generate(), userID, r.id)
			if err != nil {
				return err
			}

			r.eventManager.RecordEvent(ev)
			return nil
		}
	}
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

func (r *Room) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Room) MaxMemberCount() int {
	return r.maxMemberCount
}

func (r *Room) MemberCount() int {
	return r.memberCount
}

func (r *Room) Owner() kernel.UserID {
	return r.owner
}

func (r *Room) Number() RoomNumber {
	return r.number
}

func (r *Room) PasswordEncrypted() string {
	return r.passwordEncrypted
}

func (r *Room) Members() []kernel.UserID {
	return r.members
}

func (r *Room) GetEvents() []event.Event {
	return r.eventManager.GetEvents()
}
