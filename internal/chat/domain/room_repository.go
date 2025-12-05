//go:generate mockgen -source=room_repository.go -destination=./mocks/mock_room_repository.go -package=mocks
package domain

import (
	"context"
)

type RoomRepository interface {
	RoomCreator
}

type RoomCreator interface {
	Create(ctx context.Context, room *Room) error
}
