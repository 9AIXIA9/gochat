package domain

import "gochat/internal/shared/kernel"

type RoomID kernel.ID

func (i RoomID) String() string {
	return string(i)
}
