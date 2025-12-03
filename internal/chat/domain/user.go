package domain

import (
	"gochat/internal/shared/kernel"
)

type User struct {
	id kernel.UserID
}

func LoadUser(
	id kernel.UserID,
) *User {
	return &User{
		id: id,
	}
}

func CreateUser(
	id kernel.UserID,
) *User {
	return &User{
		id: id,
	}
}

func (u *User) ID() kernel.UserID {
	return u.id
}
