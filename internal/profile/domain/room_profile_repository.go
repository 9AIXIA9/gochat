//go:generate mockgen -source=room_profile_repository.go -destination=./mocks/mock_room_profile_repository.go -package=mocks
package domain

import (
	"context"
)

type RoomProfileRepository interface {
	RoomProfileCreator
}

type RoomProfileCreator interface {
	Create(ctx context.Context, profile *RoomProfile) error
}
