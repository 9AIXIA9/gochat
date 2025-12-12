//go:generate mockgen -source=room_profile_repository.go -destination=./mocks/mock_room_profile_repository.go -package=mocks
package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomProfileRepository interface {
	RoomProfileCreator
	RoomProfileFinder
	RoomProfileUpdater
}

type RoomProfileCreator interface {
	Create(ctx context.Context, profile *RoomProfile) error
}

type RoomProfileUpdater interface {
	Update(ctx context.Context, profile *RoomProfile) error
}

type RoomProfileFinder interface {
	FindByID(ctx context.Context, id kernel.RoomID) (*RoomProfile, error)
}
