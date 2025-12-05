//go:generate mockgen -source=roomship_repository.go -destination=./mocks/mock_roomship_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomshipRepository interface {
	RoomshipCreator
	RoomshipExisterByUserIDAndRoomID
	RoomshipsFinderByRoomID
}

type RoomshipCreator interface {
	Create(ctx context.Context, roomship *Roomship) error
}

type RoomshipExisterByUserIDAndRoomID interface {
	ExistByUserIDAndRoomID(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID) (bool, error)
}

type RoomshipsFinderByRoomID interface {
	FindsByRoomID(ctx context.Context, id kernel.RoomID) ([]*Roomship, error)
}
