package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type RoomMembersFinder interface {
	FindMembersByRoomID(ctx context.Context, roomID domain.RoomID) ([]kernel.UserID, error)
}
