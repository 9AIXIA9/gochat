//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package domain

import (
	"gochat/internal/shared/kernel"
)

type Encryptor interface {
	Encrypt(origin string) (encrypted string, err error)
}

type Comparator interface {
	Compare(encrypted string, origin string) error
}

type RoomNumberGenerator interface {
	Generate() RoomNumber
}

type RoomIDGenerator interface {
	Generate() kernel.RoomID
}

type RoomshipIDGenerator interface {
	Generate() RoomshipID
}
