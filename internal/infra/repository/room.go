package repository

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gochat/internal/utils"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) domain.RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Save(ctx context.Context, room *domain.Room) error {
	gormRoom := model.RoomFromDomain(room)
	if err := r.db.WithContext(ctx).Save(gormRoom).Error; err != nil {
		return utils.HandleDatabaseError(ctx, err)
	}
	return nil
}

func (r *RoomRepository) FindOneByNumber(ctx context.Context, roomNumber domain.RoomNumber) (*domain.Room, error) {
	var gormRoom model.Room
	if err := r.db.WithContext(ctx).Where("number = ?", roomNumber).First(&gormRoom).Error; err != nil {
		return nil, utils.HandleDatabaseError(ctx, err)
	}

	return gormRoom.ToDomain(), nil
}
