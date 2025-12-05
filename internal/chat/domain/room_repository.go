//go:generate mockgen -source=room_repository.go -destination=./mocks/mock_room_repository.go -package=mocks
package domain

import (
	"context"
)

type RoomRepository interface {
	RoomSaver
}

type RoomSaver interface {
	Save(ctx context.Context, room *Room) error
}
