package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/profile/domain"
	"gochat/internal/profile/infrastructure/persistence/model"

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
	return gormutils.TranslateError(repo.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				DoNothing: true,
			},
		).
		Create(repo.toModel(room)).Error)
}

func (repo *RoomRepository) toModel(room *domain.Room) *model.Room {
	return &model.Room{
		ID: room.ID(),
	}
}
