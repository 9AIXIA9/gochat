package repository

import (
	"errors"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gorm.io/gorm"
)

type UserRoomRepository struct {
	db *gorm.DB
}

func NewUserRoomRepository(db *gorm.DB) domain.UserRoomRepository {
	return &UserRoomRepository{db: db}
}

func (u *UserRoomRepository) Join(userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	var userRoom model.GormUserRoom
	result := u.db.Where("user_number = ? AND room_number = ?", userNumber, roomNumber).First(&userRoom)
	if result.RowsAffected > 0 {
		return true, nil // 用户已在房间中
	}

	// 创建新的用户房间关联
	newUserRoom := model.GormUserRoom{
		UserNumber: userNumber,
		RoomNumber: roomNumber,
	}

	if err := u.db.Create(&newUserRoom).Error; err != nil {
		return false, err
	}

	// 更新房间的当前用户数
	if err := u.db.Model(&model.GormRoom{}).Where("number = ?", roomNumber).
		UpdateColumn("current_users", gorm.Expr("current_users + ?", 1)).Error; err != nil {
		return false, err
	}

	return false, nil
}

func (u *UserRoomRepository) Leave(userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	result := u.db.Where("user_number = ? AND room_number = ?", userNumber, roomNumber).Delete(&model.GormUserRoom{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil // 用户不在房间中
		}
		return false, result.Error
	}

	// 更新房间的当前用户数
	if err := u.db.Model(&model.GormRoom{}).Where("number = ?", roomNumber).
		UpdateColumn("current_users", gorm.Expr("current_users - ?", 1)).Error; err != nil {
		return false, err
	}

	return true, nil
}
