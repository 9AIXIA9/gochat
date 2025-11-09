package application

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type RoomRepository interface {
	RoomIDSaver
	RoomMemberSaver
	RoomMemberDeleter
}

type RoomIDSaver interface {
	SaveID(ctx context.Context, roomID domain.RoomID) error
}

type RoomMemberSaver interface {
	SaveMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error
}

type RoomMemberDeleter interface {
	DeleteMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error
}
