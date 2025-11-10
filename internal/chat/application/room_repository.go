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
	RoomMemberFinder
	RoomMemberDeleter
}

type RoomFinder interface {
	FindByNumber(ctx context.Context, number domain.RoomNumber) (*domain.Room, error)
}

type RoomNumberSaver interface {
	SaveNumber(ctx context.Context, roomID domain.RoomID, number domain.RoomNumber) error
}

type RoomMemberSaver interface {
	SaveMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error
}

type RoomMemberFinder interface {
	FindMember(ctx context.Context, roomID domain.RoomID) ([]kernel.UserID, error)
}

type RoomMemberDeleter interface {
	DeleteMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error
}
