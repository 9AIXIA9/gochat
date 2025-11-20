package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const (
	defaultMaxMemberCount = 20
)

type Room struct {
	id                kernel.RoomID
	number            kernel.RoomNumber
	owner             kernel.UserID
	passwordEncrypted PasswordEncrypted
	members           []kernel.UserID
	maxMemberCount    int
	createdAt         time.Time
	eventManager      *event.Manager
}

type RoomOption struct {
	MaxMemberCount    int
	PasswordEncrypted PasswordEncrypted
}

func (o *RoomOption) Validate() error {
	if o.MaxMemberCount < 2 {
		return myErrors.ErrLessThanMin
	}
	return nil
}

func LoadRoom(
	id kernel.RoomID,
	owner kernel.UserID,
	number kernel.RoomNumber,
	passwordEncrypted PasswordEncrypted,
	members []kernel.UserID,
	maxMemberCount int,
	createdAt time.Time,
) *Room {
	return &Room{
		id:                id,
		owner:             owner,
		number:            number,
		passwordEncrypted: passwordEncrypted,
		members:           members,
		maxMemberCount:    maxMemberCount,
		createdAt:         createdAt,
		eventManager:      event.NewEventManager(),
	}
}

func CreateRoom(
	owner kernel.UserID,
	roomIDGenerator RoomIDGenerator,
	roomNumberGenerator RoomNumberGenerator,
	eventIDGenerator event.IDGenerator,
	options ...*RoomOption,
) (*Room, error) {
	var opt *RoomOption
	if len(options) == 0 {
		opt = &RoomOption{
			MaxMemberCount:    defaultMaxMemberCount,
			PasswordEncrypted: "",
		}
	} else {
		opt = options[0]
	}

	if err := opt.Validate(); err != nil {
		return nil, err
	}

	r := &Room{
		id:                roomIDGenerator.Generate(),
		number:            roomNumberGenerator.Generate(),
		owner:             owner,
		passwordEncrypted: opt.PasswordEncrypted,
		members:           make([]kernel.UserID, 0, 1),
		maxMemberCount:    opt.MaxMemberCount,
		createdAt:         time.Now().UTC(),
		eventManager:      event.NewEventManager(),
	}

	r.members = append(r.members, owner)

	ev, err := NewRoomCreatedEvent(eventIDGenerator.Generate(), r.id)
	if err != nil {
		return nil, err
	}
	r.eventManager.RecordEvent(ev)
	return r, nil
}

func (r *Room) AddMember(
	userID kernel.UserID,
	rawPassword Password,
	comparator Comparator,
	generator event.IDGenerator,
) error {
	if len(r.passwordEncrypted) != 0 {
		if err := comparator.Compare(r.passwordEncrypted.String(), rawPassword.String()); err != nil {
			return err
		}
	}

	if r.maxMemberCount <= r.MemberCount() {
		return myErrors.ErrExceedMaxValue
	}

	if r.IsMember(userID) {
		return nil
	}

	r.members = append(r.members, userID)

	ev, err := NewRoomJoinedEvent(generator.Generate(), userID, r.id)
	if err != nil {
		return err
	}

	r.eventManager.RecordEvent(ev)
	return nil
}

func (r *Room) DeleteMember(
	userID kernel.UserID,
	generator event.IDGenerator,
) error {
	if userID == r.owner {
		return myErrors.ErrOwnerCantLeave
	}

	for i, member := range r.members {
		if member == userID {
			r.members = append(r.members[:i], r.members[i+1:]...)

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

func (r *Room) ID() kernel.RoomID {
	return r.id
}

func (r *Room) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Room) MaxMemberCount() int {
	return r.maxMemberCount
}

func (r *Room) MemberCount() int {
	return len(r.members)
}

func (r *Room) Owner() kernel.UserID {
	return r.owner
}

func (r *Room) Number() kernel.RoomNumber {
	return r.number
}

func (r *Room) PasswordEncrypted() PasswordEncrypted {
	return r.passwordEncrypted
}

func (r *Room) Members() []kernel.UserID {
	return r.members
}

func (r *Room) GetEvents() []event.Event {
	return r.eventManager.GetEvents()
}
