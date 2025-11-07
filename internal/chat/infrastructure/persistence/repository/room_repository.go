package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ application.RoomFinder = (*RoomRepository)(nil)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (repo *RoomRepository) FindByNumber(ctx context.Context, number domain.RoomNumber) (*domain.Room, error) {
	var room model.Room
	if err := repo.db.WithContext(ctx).Preload("Members").First(&room, "number = ?", number).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}

	members := make([]kernel.UserID, 0, len(room.Members))
	for _, member := range room.Members {
		members = append(members, member.ID)
	}

	return domain.NewRoom(room.ID, room.Number, members), nil
}
