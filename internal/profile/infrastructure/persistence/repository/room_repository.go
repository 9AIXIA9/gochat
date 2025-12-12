package repository

import (
	"context"
	"gochat/internal/profile/domain"

	"gorm.io/gorm"
)

var _ domain.RoomRepository = (*RoomRepository)(nil)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (repo *RoomRepository) Save(ctx context.Context, room *domain.Room) error {
	//TODO implement me
	panic("implement me")
}
