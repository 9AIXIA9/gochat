package application

import (
	"gochat/internal/social/domain"
)

type Encryptor interface {
	Encrypt(origin string) (encrypted string, err error)
}

type Comparator interface {
	Compare(encrypted string, origin string) error
}

type RoomNumberGenerator interface {
	Generate() domain.RoomNumber
}

type RoomIDGenerator interface {
	Generate() domain.RoomID
}
