package domain

import (
	"gochat/internal/shared/kernel"
)

type Room struct {
	id kernel.RoomID
}

func CreateRoom(
	id kernel.RoomID,
) *Room {
	return &Room{id: id}
}

func (r *Room) ID() kernel.RoomID {
	return r.id
}
