package kafka

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type UserNumberSaver interface {
	SaveNumber(ctx context.Context, userID kernel.UserID, number domain.UserNumber) error
}

type RoomNumberSaver interface {
	SaveNumber(ctx context.Context, roomID domain.RoomID, number domain.RoomNumber) error
}

type RoomMemberSaver interface {
	SaveMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error
}

type RoomMemberDeleter interface {
	DeleteMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error
}
