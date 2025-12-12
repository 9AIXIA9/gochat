//go:generate mockgen -source=roomship_repository.go -destination=./mocks/mock_roomship_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomshipRepository interface {
	RoomshipSaver
	RoomshipFinderByUserIDAndRoomID
}

type RoomshipSaver interface {
	Save(ctx context.Context, roomship *Roomship) error
}

type RoomshipFinderByUserIDAndRoomID interface {
	FindByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (*Roomship, error)
}
