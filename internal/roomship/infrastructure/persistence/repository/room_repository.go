package repository

import (
	"gochat/internal/shared/event"

	"gochat/internal/roomship/domain"

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
