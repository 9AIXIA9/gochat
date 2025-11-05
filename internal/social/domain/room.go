package domain

import "gochat/internal/shared/kernel"

type RoomID kernel.ID

func (i RoomID) String() string {
	return string(i)
}

type RoomNumber kernel.Number

func (n RoomNumber) String() string {
	return string(n)
}

type Room struct {
	iD                RoomID
	number            RoomNumber
	passwordEncrypted string
}

func NewRoom(iD RoomID, number RoomNumber, passwordEncrypted string) *Room {
	return &Room{
		iD:                iD,
		number:            number,
		passwordEncrypted: passwordEncrypted,
	}
}

func (r *Room) ID() RoomID {
	return r.iD
}

func (r *Room) Number() RoomNumber {
	return r.number
}

func (r *Room) PasswordEncrypted() string {
	return r.passwordEncrypted
}
