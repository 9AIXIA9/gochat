package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomRepository interface {
	RoomSaver
	RoomFinderByNumber
	RoomFinderByID
	RoomMemberSaver
	RoomMemberDeleter
}

type RoomSaver interface {
	Save(ctx context.Context, room *Room) error
}

type RoomFinderByNumber interface {
	FindByNumber(ctx context.Context, number kernel.RoomNumber) (*Room, error)
}

type RoomFinderByID interface {
	FindByID(ctx context.Context, id kernel.RoomID) (*Room, error)
}

type RoomMemberSaver interface {
	SaveMember(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID) error
}

type RoomMemberDeleter interface {
	DeleteMember(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID) error
}
