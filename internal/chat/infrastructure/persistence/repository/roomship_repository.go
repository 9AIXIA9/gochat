package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomshipRepository = (*RoomshipRepository)(nil)

type RoomshipRepository struct {
	db *gorm.DB
}

func NewRoomshipRepository(db *gorm.DB) *RoomshipRepository {
	return &RoomshipRepository{
		db: db,
	}
}

func (repo *RoomshipRepository) Create(ctx context.Context, roomship *domain.Roomship) error {
	//TODO implement me
	panic("implement me")
}

func (repo *RoomshipRepository) ExistByUserIDAndRoomID(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *RoomshipRepository) FindsByRoomID(ctx context.Context, id kernel.RoomID) ([]*domain.Roomship, error) {
	//TODO implement me
	panic("implement me")
}
