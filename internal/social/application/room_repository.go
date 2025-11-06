package application

import (
	"context"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
)

type RoomSaver interface {
	Save(ctx context.Context, room *domain.Room) error
}

type RoomJoiner interface {
	Join(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error
}

type RoomLeaver interface {
	Leave(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error
}

type RoomFinder interface {
	FindByNumber(ctx context.Context, number domain.RoomNumber) (*domain.Room, error)
}

//type RoomMemberUpdater interface {
//	Update(ctx context.Context,)
//}
