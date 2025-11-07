package kafka

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type UserEmailSaver interface {
	SaveEmail(ctx context.Context, userID kernel.UserID, email kernel.Email) error
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
