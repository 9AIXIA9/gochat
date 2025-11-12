package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"

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
