package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomRepository interface {
	RoomFinderByID
	RoomCreator
	RoomUpdater
}

type RoomFinderByID interface {
	FindByID(ctx context.Context, id kernel.RoomID) (*Room, error)
}

type RoomCreator interface {
	Create(ctx context.Context, room *Room) error
}

type RoomUpdater interface {
	Update(ctx context.Context, room *Room) error
}
