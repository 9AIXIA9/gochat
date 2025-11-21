package domain

import (
	"gochat/internal/shared/kernel"
)

type User struct {
	id     kernel.UserID
	number kernel.UserNumber
}

func LoadUser(
	id kernel.UserID,
	number kernel.UserNumber,
) *User {
	return &User{
		id:     id,
		number: number,
	}
}

func (u *User) ID() kernel.UserID {
	return u.id
}

func (u *User) Number() kernel.UserNumber {
	return u.number
}
