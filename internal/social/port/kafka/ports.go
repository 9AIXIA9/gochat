package kafka

import (
	"context"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
)

type UserIDSaver interface {
	SaveID(ctx context.Context, id kernel.UserID) error
}

type RoomFinderByID interface {
	FindByID(ctx context.Context, id domain.RoomID) (*domain.Room, error)
}
