package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/infrastructure/persistence/model"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomshipRepository = (*RoomshipRepository)(nil)

type RoomshipRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewRoomshipRepository(db *gorm.DB, eventRepo event.Repository) *RoomshipRepository {
	return &RoomshipRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *RoomshipRepository) Create(ctx context.Context, roomship *domain.Roomship) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(repo.toModel(roomship)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		if err := repo.eventRepo.CreateUnpublishedEvents(gormutils.SetTransaction(ctx, tx), roomship.GetEvents()); err != nil {
			return err
		}
		return nil
	})
}

func (repo *RoomshipRepository) FindByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (*domain.Roomship, error) {
	var roomship model.Roomship
	if err := repo.db.WithContext(ctx).Where("user_id = ? AND room_id = ?", userID, roomID).First(&roomship).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&roomship), nil
}

func (repo *RoomshipRepository) ExistByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (bool, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&model.Roomship{}).Where("user_id = ? AND room_id = ?", userID, roomID).Count(&count).Error; err != nil {
		return false, gormutils.TranslateError(err)
	}
	return count > 0, nil
}

func (repo *RoomshipRepository) toModel(roomship *domain.Roomship) *model.Roomship {
	return &model.Roomship{
		ID:        roomship.ID(),
		RoomID:    roomship.RoomID(),
		UserID:    roomship.UserID(),
		Role:      roomship.Role(),
		CreatedAt: roomship.CreatedAt(),
	}
}

func (repo *RoomshipRepository) toDomain(roomship *model.Roomship) *domain.Roomship {
	return domain.LoadRoomship(roomship.ID, roomship.RoomID, roomship.UserID, roomship.Role, roomship.CreatedAt)
}
