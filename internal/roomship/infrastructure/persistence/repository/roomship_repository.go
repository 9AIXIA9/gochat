package repository

import (
	"context"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomshipRepository = (*RoomshipRepository)(nil)

type RoomshipRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewRoomshipRepository(db *gorm.DB, eventRepo event.Repository) *RoomshipRepository {
	return &RoomshipRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *RoomshipRepository) Create(ctx context.Context, roomship *domain.Roomship) error {
	//TODO implement me
	panic("implement me")
}

func (repo *RoomshipRepository) FindByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (*domain.Roomship, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *RoomshipRepository) ExistByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (bool, error) {
	//TODO implement me
	panic("implement me")
}
