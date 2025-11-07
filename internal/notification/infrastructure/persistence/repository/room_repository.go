package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/notification/port/kafka"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ kafka.RoomIDSaver = (*RoomRepository)(nil)
var _ kafka.RoomMemberSaver = (*RoomRepository)(nil)
var _ kafka.RoomMemberDeleter = (*RoomRepository)(nil)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (repo *RoomRepository) SaveID(ctx context.Context, roomID domain.RoomID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&model.Room{
		ID: roomID,
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
