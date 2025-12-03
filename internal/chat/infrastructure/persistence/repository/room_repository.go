package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
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

func (repo *RoomRepository) Create(ctx context.Context, room *domain.Room) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(repo.toModel(room)).Error)
}

func (repo *RoomRepository) Update(ctx context.Context, room *domain.Room) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		modelRoom := repo.toModel(room)

		// 使用 Replace 来更新多对多关联的成员
		// 删除不再存在的关联并添加新的关联
		err := tx.Model(&modelRoom).Association("Members").Replace(modelRoom.Members)
		if err != nil {
			return gormutils.TranslateError(err)
		}

		// 更新 Room 表本身的字段
		if err := tx.Updates(modelRoom).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *RoomRepository) FindByID(ctx context.Context, id kernel.RoomID) (*domain.Room, error) {
	var room model.Room
	if err := repo.db.WithContext(ctx).Preload("Members").First(&room, "id = ?", id).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&room), nil
}

func (repo *RoomRepository) toModel(room *domain.Room) *model.Room {
	members := make([]*model.User, 0, len(room.Members()))
	for _, memberID := range room.Members() {
		members = append(members, &model.User{
			ID: memberID,
		})
	}
	return &model.Room{
		ID:      room.ID(),
		Members: members,
	}
}

func (repo *RoomRepository) toDomain(room *model.Room) *domain.Room {
	members := make([]kernel.UserID, 0, len(room.Members))
	for _, member := range room.Members {
		members = append(members, member.ID)
	}

	return domain.LoadRoom(room.ID, members)
}
