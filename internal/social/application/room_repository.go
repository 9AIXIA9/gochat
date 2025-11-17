package application

import (
	"context"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
)

type RoomRepository interface {
	RoomSaver
	RoomFinderByNumber
	RoomFinderByID
	RoomMemberSaver
	RoomMemberDeleter
}

type RoomSaver interface {
	Save(ctx context.Context, room *domain.Room) error
}

type RoomFinderByNumber interface {
	FindByNumber(ctx context.Context, number kernel.RoomNumber) (*domain.Room, error)
}

type RoomFinderByID interface {
	FindByID(ctx context.Context, id kernel.RoomID) (*domain.Room, error)
}

type RoomMemberSaver interface {
	SaveMember(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID) error
}

type RoomMemberDeleter interface {
	DeleteMember(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID) error
}
