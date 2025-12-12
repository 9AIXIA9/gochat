package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/profile/domain"
	"gochat/internal/profile/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomProfileRepository = (*RoomProfileRepository)(nil)

type RoomProfileRepository struct {
	db *gorm.DB
}

func NewRoomProfileRepository(db *gorm.DB) *RoomProfileRepository {
	return &RoomProfileRepository{db: db}
}

func (repo *RoomProfileRepository) Create(ctx context.Context, profile *domain.RoomProfile) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(repo.toModel(profile)).Error)
}

func (repo *RoomProfileRepository) Update(ctx context.Context, profile *domain.RoomProfile) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Updates(repo.toModel(profile)).Error)
}

func (repo *RoomProfileRepository) FindByID(ctx context.Context, id kernel.RoomID) (*domain.RoomProfile, error) {
	var modelRoomProfile model.RoomProfile
	if err := repo.db.WithContext(ctx).First(&modelRoomProfile, "id = ?", id).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}

	return repo.toDomain(&modelRoomProfile), nil
}

func (repo *RoomProfileRepository) toModel(profile *domain.RoomProfile) *model.RoomProfile {
	return &model.RoomProfile{
		ID:           profile.ID(),
		Name:         profile.Name(),
		Introduction: profile.Introduction(),
		CreatedAt:    profile.CreatedAt(),
	}
}

func (repo *RoomProfileRepository) toDomain(profile *model.RoomProfile) *domain.RoomProfile {
	return domain.LoadRoomProfile(
		profile.ID,
		profile.Name,
		profile.Introduction,
		profile.CreatedAt,
	)
}
