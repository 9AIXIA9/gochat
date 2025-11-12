package application

import (
	"context"
	"gochat/internal/notification/domain"
)

type RoomRepository interface {
	RoomIDSaver
}

type RoomIDSaver interface {
	SaveID(ctx context.Context, roomID domain.RoomID) error
}
