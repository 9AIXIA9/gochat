package application

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type RoomRepository interface {
	RoomFinder
	RoomNumberSaver
	RoomMemberSaver
	RoomMemberDeleter
}

type RoomFinder interface {
	FindByNumber(ctx context.Context, number kernel.RoomNumber) (*domain.Room, error)
}

type RoomNumberSaver interface {
	SaveNumber(ctx context.Context, roomID kernel.RoomID, number kernel.RoomNumber) error
}

type RoomMemberSaver interface {
	SaveMember(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID) error
}

type RoomMemberDeleter interface {
	DeleteMember(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID) error
}
