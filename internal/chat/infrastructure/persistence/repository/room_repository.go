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

var _ application.RoomRepository = (*RoomRepository)(nil)

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

func (repo *RoomRepository) SaveNumber(ctx context.Context, roomID domain.RoomID, number domain.RoomNumber) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&model.Room{
		ID:     roomID,
		Number: number,
	}).Error)
}

func (repo *RoomRepository) SaveMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	room := model.Room{ID: roomID}
	user := model.User{ID: userID}

	return gormutils.TranslateError(
		repo.db.WithContext(ctx).
			Model(&room).
			Association("Members").
			Append(&user),
	)
}

func (repo *RoomRepository) FindMember(ctx context.Context, roomID domain.RoomID) ([]kernel.UserID, error) {
	var members []model.User
	if err := repo.db.WithContext(ctx).
		Model(&model.Room{ID: roomID}).
		Association("Members").
		Find(&members); err != nil {
		return nil, gormutils.TranslateError(err)
	}

	ids := make([]kernel.UserID, 0, len(members))
	for _, m := range members {
		ids = append(ids, m.ID)
	}
	return ids, nil
}

func (repo *RoomRepository) DeleteMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	room := model.Room{ID: roomID}
	user := model.User{ID: userID}

	return gormutils.TranslateError(
		repo.db.WithContext(ctx).
			Model(&room).
			Association("Members").
			Delete(&user),
	)
}
