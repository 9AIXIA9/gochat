package repository

import (
	"context"
	"errors"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
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
	return gormutils.TranslateError(repo.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
		DoNothing: true,
	}).Create(&model.Room{
		ID:     roomID,
		Number: number,
	}).Error)
}

func (repo *RoomRepository) SaveMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	if err := gormutils.TranslateError(
		repo.db.WithContext(ctx).
			Model(&model.Room{ID: roomID}).
			Association("Members").
			Append(&model.User{ID: userID}),
	); err != nil {
		if errors.Is(err, myErrors.ErrDuplicatedKey) {
			return nil
		}
		return err
	}
	return nil
}

func (repo *RoomRepository) DeleteMember(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	if err := gormutils.TranslateError(
		repo.db.WithContext(ctx).
			Model(&model.Room{ID: roomID}).
			Association("Members").
			Delete(&model.User{ID: userID}),
	); err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil
		}
		return err
	}
	return nil
}
