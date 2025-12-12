//go:generate mockgen -source=roomship_repository.go -destination=./mocks/mock_roomship_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomshipRepository interface {
	RoomshipCreator
	RoomshipFinderByID
	RoomshipsFinderByRoomID
	RoomshipsFinderByRoomIDAndRole
	RoomshipFinderByUserIDAndRoomID
	RoomshipExisterByUserIDAndRoomID
	RoomshipsFinderByUserID
	RoomshipDeleterByUserIDAndRoomID
}

type RoomshipCreator interface {
	Create(ctx context.Context, roomship *Roomship) error
}

type RoomshipFinderByID interface {
	FindByID(ctx context.Context, id RoomshipID) (*Roomship, error)
}

type RoomshipsFinderByRoomID interface {
	FindsByRoomID(ctx context.Context, id kernel.RoomID) ([]*Roomship, error)
}

type RoomshipsFinderByRoomIDAndRole interface {
	FindsByRoomIDAndRole(ctx context.Context, roomID kernel.RoomID, role RoomshipRole) ([]*Roomship, error)
}

type RoomshipFinderByUserIDAndRoomID interface {
	FindByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (*Roomship, error)
}

type RoomshipExisterByUserIDAndRoomID interface {
	ExistByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (bool, error)
}

type RoomshipsFinderByUserID interface {
	FindsByUserID(ctx context.Context, userID kernel.UserID, limit int, baseID RoomshipID) ([]*Roomship, error)
}

type RoomshipDeleterByUserIDAndRoomID interface {
	DeleteByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) error
}
