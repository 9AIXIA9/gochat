package domain

import (
	"context"
	"gochat/internal/shared/kernel"
)

//TODO 存储层也要进行重构

type RoomRepository interface {
	RoomFinder
	RoomNumberSaver
	RoomMemberSaver
	RoomMemberDeleter
	RoomMembersFinder
}

type RoomFinder interface {
	FindByNumber(ctx context.Context, number kernel.RoomNumber) (*Room, error)
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

type RoomMembersFinder interface {
	FindMembers(ctx context.Context, roomID kernel.RoomID) ([]kernel.UserID, error)
}
