package domain

import (
	"gochat/internal/shared/kernel"
)

type Room struct {
	id     kernel.RoomID
	number kernel.RoomNumber

	members []kernel.UserID
}

func LoadRoom(
	id kernel.RoomID,
	number kernel.RoomNumber,
	members []kernel.UserID,
) *Room {
	return &Room{
		id:      id,
		number:  number,
		members: members,
	}
}

func (r *Room) ID() kernel.RoomID {
	return r.id
}

func (r *Room) Number() kernel.RoomNumber {
	return r.number
}

func (r *Room) Members() []kernel.UserID {
	return r.members
}

func (r *Room) IsMember(userID kernel.UserID) bool {
	for _, memberID := range r.members {
		if memberID == userID {
			return true
		}
	}
	return false
}
