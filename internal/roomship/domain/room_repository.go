//go:generate mockgen -source=room_repository.go -destination=./mocks/mock_room_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomRepository interface {
	RoomCreator
	RoomUpdater
	RoomFinderByNumber
	RoomFinderByID
}

type RoomCreator interface {
	Create(ctx context.Context, room *Room) error
}

type RoomUpdater interface {
	Update(ctx context.Context, room *Room) error
}

type RoomFinderByNumber interface {
	FindByNumber(ctx context.Context, number kernel.RoomNumber) (*Room, error)
}

type RoomFinderByID interface {
	FindByID(ctx context.Context, id kernel.RoomID) (*Room, error)
}
