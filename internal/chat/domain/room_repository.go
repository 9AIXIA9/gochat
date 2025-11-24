package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomRepository interface {
	RoomFinderByID
	RoomFinderByNumber
	RoomCreator
	RoomUpdater
}

type RoomFinderByID interface {
	FindByID(ctx context.Context, id kernel.RoomID) (*Room, error)
}

type RoomFinderByNumber interface {
	FindByNumber(ctx context.Context, number kernel.RoomNumber) (*Room, error)
}

type RoomCreator interface {
	Create(ctx context.Context, room *Room) error
}

type RoomUpdater interface {
	Update(ctx context.Context, room *Room) error
}
