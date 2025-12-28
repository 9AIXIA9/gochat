package domain

import "gochat/internal/shared/kernel"

type RoomshipID string

func (id RoomshipID) String() string {
	return string(id)
}

type Roomship struct {
	id     RoomshipID
	userID kernel.UserID
	roomID kernel.RoomID
}

func LoadRoomship(
	id RoomshipID,
	userID kernel.UserID,
	roomID kernel.RoomID,
) *Roomship {
	return &Roomship{
		id:     id,
		userID: userID,
		roomID: roomID,
	}
}

func (f *Roomship) ID() RoomshipID {
	return f.id
}

func (f *Roomship) UserID() kernel.UserID {
	return f.userID
}

func (f *Roomship) RoomID() kernel.RoomID {
	return f.roomID
}
