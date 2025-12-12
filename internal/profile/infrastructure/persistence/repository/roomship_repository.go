package repository

import (
	"context"
	"gochat/internal/profile/domain"
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

func (repo *RoomshipRepository) Save(ctx context.Context, roomship *domain.Roomship) error {
	//TODO implement me
	panic("implement me")
}

func (repo *RoomshipRepository) FindByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (*domain.Roomship, error) {
	//TODO implement me
	panic("implement me")
}
