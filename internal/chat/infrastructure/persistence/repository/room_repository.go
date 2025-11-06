package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"

	"gorm.io/gorm"
)

var _ application.RoomFinder = (*RoomRepository)(nil)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r RoomRepository) FindByNumber(ctx context.Context, number domain.RoomNumber) (*domain.Room, error) {
	//TODO implement me
	panic("implement me")
}
