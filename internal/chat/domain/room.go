package domain

import (
	"gochat/internal/shared/kernel"
)

type Room struct {
	id kernel.RoomID

	members []kernel.UserID
}

func CreateRoom(
	id kernel.RoomID,
	ownerID kernel.UserID,
) *Room {
	return &Room{
		id:      id,
		members: []kernel.UserID{ownerID},
	}
}

func LoadRoom(
	id kernel.RoomID,
	members []kernel.UserID,
) *Room {
	return &Room{
		id:      id,
		members: members,
	}
}

func (r *Room) AddMember(
	userID kernel.UserID,
) {
	if r.IsMember(userID) {
		return
	}
	r.members = append(r.members, userID)
}

func (r *Room) DeleteMember(
	userID kernel.UserID,
) {
	for i, member := range r.members {
		if member == userID {
			r.members = append(r.members[:i], r.members[i+1:]...)
			return
		}
	}
}

func (r *Room) ID() kernel.RoomID {
	return r.id
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
