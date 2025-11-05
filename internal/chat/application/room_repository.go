package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type RoomMembersFinder interface {
	FindByRoomID(ctx context.Context, roomID domain.RoomID) ([]kernel.UserID, error)
}
