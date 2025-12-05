package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomshipRepository = (*RoomshipRepository)(nil)

type RoomshipRepository struct {
	db *gorm.DB
}

func NewRoomshipRepository(db *gorm.DB) *RoomshipRepository {
	return &RoomshipRepository{
		db: db,
	}
}

func (repo *RoomshipRepository) Create(ctx context.Context, roomship *domain.Roomship) error {
	if err := repo.db.WithContext(ctx).Create(repo.toModel(roomship)).Error; err != nil {
		return gormutils.TranslateError(err)
	}
	return nil
}

func (repo *RoomshipRepository) ExistByUserIDAndRoomID(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&model.Roomship{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Count(&count).Error
	if err != nil {
		return false, gormutils.TranslateError(err)
	}

	return count > 0, nil
}

func (repo *RoomshipRepository) FindsByRoomID(ctx context.Context, id kernel.RoomID) ([]*domain.Roomship, error) {
	var roomships []model.Roomship
	err := repo.db.WithContext(ctx).
		Model(&model.Roomship{}).
		Where("room_id = ?", id).
		Find(&roomships).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}

	return repo.toDomains(roomships), nil
}

func (repo *RoomshipRepository) toModel(roomship *domain.Roomship) *model.Roomship {
	return &model.Roomship{
		ID:     roomship.ID(),
		UserID: roomship.UserID(),
		RoomID: roomship.RoomID(),
	}
}

func (repo *RoomshipRepository) toDomain(roomship *model.Roomship) *domain.Roomship {
	return domain.LoadRoomship(roomship.ID, roomship.UserID, roomship.RoomID)
}

func (repo *RoomshipRepository) toDomains(roomships []model.Roomship) []*domain.Roomship {
	domains := make([]*domain.Roomship, len(roomships))
	for i, roomship := range roomships {
		domains[i] = repo.toDomain(&roomship)
	}
	return domains
}
