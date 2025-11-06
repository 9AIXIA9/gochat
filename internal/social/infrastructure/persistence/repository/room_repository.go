package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
	"gochat/internal/social/domain"
	"gochat/internal/social/infrastructure/persistence/model"

	"gorm.io/gorm"
)

var _ application.RoomFinder = (*RoomRepository)(nil)
var _ application.RoomSaver = (*RoomRepository)(nil)
var _ application.RoomJoiner = (*RoomRepository)(nil)
var _ application.RoomLeaver = (*RoomRepository)(nil)

type RoomRepository struct {
	innerRepository *gormutils.Repository[model.Room, domain.Room]
	db              *gorm.DB
}

func NewRoomRepository(db *gorm.DB, converter gormutils.GenericModelConverter[*model.Room, *domain.Room]) *RoomRepository {
	return &RoomRepository{innerRepository: gormutils.NewRepository(db, converter), db: db}
}

func (repo *RoomRepository) Save(ctx context.Context, room *domain.Room) error {
	return repo.innerRepository.Save(ctx, room)
}

func (repo *RoomRepository) FindByNumber(ctx context.Context, number domain.RoomNumber) (*domain.Room, error) {
	room, err := repo.innerRepository.Find(ctx, gormutils.Where("number = ?", number))
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (repo *RoomRepository) Join(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model.RoomMember{
			RoomID: roomID,
			Member: userID,
		}).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Room{}).Where("id = ?", roomID).UpdateColumn("member_count", gorm.Expr("member_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	}))
}

func (repo *RoomRepository) Leave(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("room_id = ? AND member = ?", roomID, userID).Delete(&model.RoomMember{})
		if err := result.Error; err != nil {
			return gormutils.TranslateError(err)
		}

		if result.RowsAffected == 0 {
			return nil
		}

		if err := tx.Model(&model.Room{}).Where("id = ?", roomID).UpdateColumn("member_count", gorm.Expr("member_count - ?", 1)).Error; err != nil {
			return err
		}
		return nil
	}))
}
