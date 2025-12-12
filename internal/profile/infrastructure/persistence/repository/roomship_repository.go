package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/profile/domain"
	"gochat/internal/profile/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (repo *RoomshipRepository) Save(ctx context.Context, roomship *domain.Roomship) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				DoNothing: true,
			},
		).
		Create(repo.toModel(roomship)).Error)
}

func (repo *RoomshipRepository) FindByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID) (*domain.Roomship, error) {
	var roomshipModel model.Roomship
	err := gormutils.TranslateError(repo.db.WithContext(ctx).
		Where("user_id = ? AND room_id = ?", userID, roomID).
		First(&roomshipModel).Error)
	if err != nil {
		return nil, err
	}

	return repo.toDomain(&roomshipModel), nil
}

func (repo *RoomshipRepository) toModel(roomship *domain.Roomship) *model.Roomship {
	return &model.Roomship{
		ID:     roomship.ID(),
		Role:   roomship.Role(),
		RoomID: roomship.RoomID(),
		UserID: roomship.UserID(),
	}
}

func (repo *RoomshipRepository) toDomain(roomship *model.Roomship) *domain.Roomship {
	return domain.LoadRoomship(roomship.ID, roomship.UserID, roomship.RoomID, roomship.Role)
}
