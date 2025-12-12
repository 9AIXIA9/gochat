package domain

import "gochat/internal/shared/kernel"

type RoomshipID kernel.ID

func (id RoomshipID) String() string {
	return string(id)
}

type RoomshipRole string

const (
	OwnerRole  RoomshipRole = "owner"
	MemberRole RoomshipRole = "member"
)

func (r RoomshipRole) String() string {
	return string(r)
}

type Roomship struct {
	id     RoomshipID
	userID kernel.UserID
	roomID kernel.RoomID
	role   RoomshipRole
}

func LoadRoomship(
	id RoomshipID,
	userID kernel.UserID,
	roomID kernel.RoomID,
	role RoomshipRole,
) *Roomship {
	return &Roomship{
		id:     id,
		userID: userID,
		roomID: roomID,
		role:   role,
	}
}

func (r *Roomship) ID() RoomshipID {
	return r.id
}

func (r *Roomship) UserID() kernel.UserID {
	return r.userID
}

func (r *Roomship) RoomID() kernel.RoomID {
	return r.roomID
}

func (r *Roomship) Role() RoomshipRole {
	return r.role
}

func (r *Roomship) IsOwner() bool {
	return r.role == OwnerRole
}
