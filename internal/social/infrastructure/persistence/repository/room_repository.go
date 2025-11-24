package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/event"

	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
	"gochat/internal/social/infrastructure/persistence/model"

	"gorm.io/gorm"
)

var _ domain.RoomRepository = (*RoomRepository)(nil)

type RoomRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewRoomRepository(db *gorm.DB, eventRepo event.Repository) *RoomRepository {
	return &RoomRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *RoomRepository) Create(ctx context.Context, room *domain.Room) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := repo.db.Create(repo.toModel(room)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		//TODO 传递事务
		txCtx := context.WithValue(ctx, gormutils.UnitOfWorkKey, tx)

		if err := repo.eventRepo.CreateUnpublishedEvents(txCtx, room.GetEvents()); err != nil {
			return gormutils.TranslateError(err)
		}
		return nil
	})
}

func (repo *RoomRepository) Update(ctx context.Context, room *domain.Room) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		modelRoom := repo.toModel(room)

		// 使用 Replace 来更新多对多关联的成员
		// 删除不再存在的关联并添加新的关联
		if err := tx.Model(&modelRoom).Association("Members").Replace(modelRoom.Members); err != nil {
			return gormutils.TranslateError(err)
		}

		// 更新 Room 表本身的字段
		if err := tx.Updates(modelRoom).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		txCtx := context.WithValue(ctx, gormutils.UnitOfWorkKey, tx)

		if err := repo.eventRepo.CreateUnpublishedEvents(txCtx, room.GetEvents()); err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *RoomRepository) FindByNumber(ctx context.Context, number kernel.RoomNumber) (*domain.Room, error) {
	var m model.Room
	if err := repo.db.WithContext(ctx).
		Preload("Members").
		Where("number = ?", number).
		First(&m).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&m), nil
}

func (repo *RoomRepository) FindByID(ctx context.Context, id kernel.RoomID) (*domain.Room, error) {
	var m model.Room
	if err := repo.db.WithContext(ctx).
		Preload("Members").
		Where("id = ?", id).
		First(&m).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&m), nil
}

func (repo *RoomRepository) toModel(room *domain.Room) *model.Room {
	modelMembers := make([]*model.User, 0, len(room.Members()))
	for _, member := range room.Members() {
		modelMembers = append(modelMembers, &model.User{
			ID: member,
		})
	}
	return &model.Room{
		ID:                room.ID(),
		OwnerID:           room.OwnerID(),
		Number:            room.Number(),
		PasswordEncrypted: room.PasswordEncrypted(),
		MaxMemberCount:    room.MaxMemberCount(),
		CreatedAt:         room.CreatedAt(),
		Owner: &model.User{
			ID: room.OwnerID(),
		},
		Members: modelMembers,
	}
}

func (repo *RoomRepository) toDomain(m *model.Room) *domain.Room {
	memberIDs := make([]kernel.UserID, 0, len(m.Members))
	for _, u := range m.Members {
		memberIDs = append(memberIDs, u.ID)
	}
	return domain.LoadRoom(
		m.ID,
		m.OwnerID,
		m.Number,
		m.PasswordEncrypted,
		memberIDs,
		m.MaxMemberCount,
		m.CreatedAt,
	)
}
