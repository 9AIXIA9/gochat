package repository

import (
	"context"

	"gochat/internal/domain"
	"gochat/internal/model"
	"gochat/internal/types"
	"gochat/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ domain.RoomRepository = (*RoomRepository)(nil)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) DoJoinRoomUOW(ctx context.Context, uow domain.JoinRoomUnitOfWork) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return uow(NewRoomRepository(tx))
	})
}

func (r *RoomRepository) FindRoom(ctx context.Context, number domain.RoomNumber) (*domain.Room, error) {
	var gormRoom model.Room
	if err := r.db.WithContext(ctx).Where("number = ?", number).First(&gormRoom).Error; err != nil {
		return nil, utils.HandleDatabaseError(ctx, err)
	}
	return gormRoom.ToDomain(), nil
}

func (r *RoomRepository) SaveRoom(ctx context.Context, room *domain.Room) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gormRoom := model.RoomFromDomain(room)
		if err := tx.Save(gormRoom).Error; err != nil {
			return utils.HandleDatabaseError(ctx, err)
		}

		if err := tx.Create(model.NewUserRoom(room.Owner(), room.Number())).Error; err != nil {
			return utils.HandleDatabaseError(ctx, err)
		}
		if err := tx.Model(&model.Room{}).Where("number = ?", room.Number()).
			UpdateColumn("current_users", gorm.Expr("current_users + ?", 1)).Error; err != nil {
			return utils.HandleDatabaseError(ctx, err)
		}

		return nil
	})
}

func (r *RoomRepository) JoinRoom(ctx context.Context, userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 尝试插入关系，避免重复插入导致重复计数
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).
			Create(model.NewUserRoom(userNumber, roomNumber))
		if err := res.Error; err != nil {
			return utils.HandleDatabaseError(ctx, err)
		}
		// 已经在房间则无需计数
		if res.RowsAffected == 0 {
			return nil
		}

		// 原子自增，不做读取，避免锁竞争
		if err := tx.Model(&model.Room{}).
			Where("number = ?", roomNumber).
			UpdateColumn("current_users", gorm.Expr("current_users + ?", 1)).Error; err != nil {
			return utils.HandleDatabaseError(ctx, err)
		}
		return nil
	})
}

func (r *RoomRepository) LeaveRoom(ctx context.Context, userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 仅在删除关系成功时再递减
		res := tx.Where("user_number = ? AND room_number = ?", userNumber, roomNumber).
			Delete(&model.UserRoom{})
		if err := res.Error; err != nil {
			return utils.HandleDatabaseError(ctx, err)
		}
		if res.RowsAffected == 0 {
			return types.ErrNotFound
		}

		// 原子自减，防止负数（可选 WHERE current_users > 0）
		if err := tx.Model(&model.Room{}).
			Where("number = ?", roomNumber).
			UpdateColumn("current_users", gorm.Expr("current_users - ?", 1)).Error; err != nil {
			return utils.HandleDatabaseError(ctx, err)
		}
		return nil
	})
}

func (r *RoomRepository) FindRoomMembers(ctx context.Context, number domain.RoomNumber) ([]domain.UserNumber, error) {
	var userRooms []model.UserRoom
	if err := utils.HandleDatabaseError(ctx, r.db.WithContext(ctx).
		Where("room_number = ?", number).Find(&userRooms).Error); err != nil {
		return nil, err
	}

	if len(userRooms) == 0 {
		return nil, nil
	}

	userNumbers := make([]domain.UserNumber, 0, len(userRooms))
	for _, ur := range userRooms {
		userNumbers = append(userNumbers, ur.UserNumber)
	}
	return userNumbers, nil
}
