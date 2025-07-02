package repository

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gochat/internal/types"
	"gochat/internal/utils"
	"gorm.io/gorm"
)

type UserRoomRepository struct {
	db *gorm.DB
}

func NewUserRoomRepository(db *gorm.DB) domain.UserRoomRepository {
	return &UserRoomRepository{db: db}
}

func (u *UserRoomRepository) Save(ctx context.Context, userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	if err := u.db.WithContext(ctx).Save(model.NewUserRoom(userNumber, roomNumber)).Error; err != nil {
		return utils.CheckDuplicateKeyError(err)
	}

	// 更新房间的当前用户数
	if err := u.db.WithContext(ctx).Model(&model.Room{}).Where("number = ?", roomNumber).
		UpdateColumn("current_users", gorm.Expr("current_users + ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRoomRepository) Delete(ctx context.Context, userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	result := u.db.WithContext(ctx).Where("user_number = ? AND room_number = ?", userNumber, roomNumber).Delete(&model.UserRoom{})
	if err := result.Error; err != nil {
		return utils.CheckNotFoundError(err)
	}

	// 如果没有删除任何记录
	if result.RowsAffected == 0 {
		return types.ErrNotFound
	}

	// 更新房间的当前用户数
	if err := u.db.WithContext(ctx).Model(&model.Room{}).Where("number = ?", roomNumber).
		UpdateColumn("current_users", gorm.Expr("current_users - ?", 1)).Error; err != nil {
		return err
	}

	return nil
}
