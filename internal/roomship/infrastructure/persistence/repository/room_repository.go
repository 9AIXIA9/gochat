package repository

import (
	"context"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomRepository = (*RoomRepository)(nil)

type RoomRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewRoomRepository(db *gorm.DB, eventRepo event.Repository) *RoomRepository {
	return &RoomRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *RoomRepository) Create(ctx context.Context, room *domain.Room) error {
	//TODO implement me
	panic("implement me")
}

func (repo *RoomRepository) FindByID(ctx context.Context, roomID kernel.RoomID) (*domain.Room, error) {
	//TODO implement me
	panic("implement me")
}
