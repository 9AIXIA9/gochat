package repository

import (
	"context"
	"errors"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ application.RoomRepository = (*RoomRepository)(nil)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (repo *RoomRepository) SaveID(ctx context.Context, roomID domain.RoomID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
			DoNothing: true,
		}).Create(&model.Room{
		ID: roomID,
	}).Error)
}

func (repo *RoomRepository) SaveMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	room := model.Room{ID: roomID}
	user := model.User{ID: userID}

	if err := gormutils.TranslateError(
		repo.db.WithContext(ctx).
			Model(&room).
			Association("Members").
			Append(&user),
	); err != nil {
		if errors.Is(err, myErrors.ErrDuplicatedKey) {
			return nil
		}
		return err
	}
	return nil
}

func (repo *RoomRepository) DeleteMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	room := model.Room{ID: roomID}
	user := model.User{ID: userID}

	if err := gormutils.TranslateError(
		repo.db.WithContext(ctx).
			Model(&room).
			Association("Members").
			Delete(&user),
	); err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil
		}
		return err
	}
	return nil
}
