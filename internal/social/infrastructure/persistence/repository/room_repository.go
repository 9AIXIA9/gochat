package repository

import (
	"context"
	"errors"
	gormutils "gochat/internal/infrastructure/gorm"
	myErrors "gochat/internal/shared/errors"

	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
	"gochat/internal/social/domain"
	"gochat/internal/social/infrastructure/persistence/model"

	"gorm.io/gorm"
)

var _ application.RoomRepository = (*RoomRepository)(nil)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (repo *RoomRepository) Save(ctx context.Context, room *domain.Room) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		modelRoom := toModelRoom(room)
		if err := tx.Create(modelRoom).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		modelMembers := make([]*model.User, 0, len(room.Members()))
		for _, member := range room.Members() {
			modelMembers = append(modelMembers, &model.User{
				ID: member,
			})
		}

		// Replace 将以传入集合为准，原子同步
		if err := tx.Model(modelRoom).Association("Members").Replace(&modelMembers); err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *RoomRepository) FindByNumber(ctx context.Context, number domain.RoomNumber) (*domain.Room, error) {
	var m model.Room
	err := repo.db.WithContext(ctx).
		Preload("Members").
		Where("number = ?", number).
		First(&m).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return toDomainRoom(&m), nil
}

func (repo *RoomRepository) FindByID(ctx context.Context, id domain.RoomID) (*domain.Room, error) {
	var m model.Room
	err := repo.db.WithContext(ctx).
		Preload("Members").
		First(&m, "id = ?", id.String()).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return toDomainRoom(&m), nil
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

func toModelRoom(room *domain.Room) *model.Room {
	return &model.Room{
		ID:                room.ID(),
		OwnerID:           room.Owner(),
		Number:            room.Number(),
		PasswordEncrypted: room.PasswordEncrypted(),
		MaxMemberCount:    room.MaxMemberCount(),
	}
}

func toDomainRoom(m *model.Room) *domain.Room {
	memberIDs := make([]kernel.UserID, 0, len(m.Members))
	for _, u := range m.Members {
		memberIDs = append(memberIDs, u.ID)
	}
	return domain.NewRoom(
		m.ID,
		m.OwnerID,
		m.Number,
		m.PasswordEncrypted,
		memberIDs,
		m.MaxMemberCount,
		m.CreatedAt,
	)
}
