package repository

import (
	"context"
	"errors"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) domain.RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, room *domain.Room) (bool, error) {
	var existingRoom model.Room
	result := r.db.WithContext(ctx).Where("number = ?", room.Number).First(&existingRoom)
	if result.RowsAffected > 0 {
		return true, nil // 房间已存在
	}

	gormRoom := model.RoomFromDomain(room)
	if err := r.db.WithContext(ctx).Create(gormRoom).Error; err != nil {
		return false, err
	}

	return false, nil
}

func (r *RoomRepository) QueryByRoomNumber(ctx context.Context, roomNumber domain.RoomNumber) (*domain.Room, error) {
	var gormRoom model.Room
	result := r.db.WithContext(ctx).Where("number = ?", roomNumber).First(&gormRoom)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 房间不存在
		}
		return nil, result.Error
	}

	return gormRoom.ToDomain(), nil
}
