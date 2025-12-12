package domain

import "gochat/internal/shared/kernel"

type RoomshipID kernel.ID

func (id RoomshipID) String() string {
	return string(id)
}

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

func (r Role) String() string {
	return string(r)
}

type Roomship struct {
	id     RoomshipID
	userID kernel.UserID
	roomID kernel.RoomID
	role   Role
}

func LoadRoomship(
	id RoomshipID,
	userID kernel.UserID,
	roomID kernel.RoomID,
	role Role,
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

func (r *Roomship) Role() Role {
	return r.role
}
