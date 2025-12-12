//go:generate mockgen -source=room_repository.go -destination=./mocks/mock_room_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomRepository interface {
	RoomCreator
	RoomFinderByID
	RoomDeleterByID
}

type RoomCreator interface {
	Create(ctx context.Context, room *Room) error
}

type RoomFinderByID interface {
	FindByID(ctx context.Context, roomID kernel.RoomID) (*Room, error)
}

type RoomDeleterByID interface {
	DeleteByID(ctx context.Context, roomID kernel.RoomID) error
}
