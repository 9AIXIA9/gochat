package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomshipRepository interface {
	RoomshipFinderByUserIDAndRoomID
	RoomshipExisterByUserIDAndRoomID
}

type RoomshipCreator interface {
	Create(ctx context.Context, roomship *Roomship) error
}

type RoomshipFinderByUserIDAndRoomID interface {
	FindByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (*Roomship, error)
}

type RoomshipExisterByUserIDAndRoomID interface {
	ExistByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (bool, error)
}
