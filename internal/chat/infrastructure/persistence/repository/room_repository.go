package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ domain.RoomRepository = (*RoomRepository)(nil)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (repo *RoomRepository) Save(ctx context.Context, room *domain.Room) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
			DoNothing: true,
		},
	).Create(&model.Room{
		ID: room.ID(),
	}).Error)
}
