package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/profile/domain"
	"gochat/internal/profile/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.UserProfileRepository = (*UserProfileRepository)(nil)

type UserProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (repo *UserProfileRepository) Create(ctx context.Context, profile *domain.UserProfile) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(repo.toModel(profile)).Error)
}

func (repo *UserProfileRepository) Update(ctx context.Context, profile *domain.UserProfile) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Updates(repo.toModel(profile)).Error)
}

func (repo *UserProfileRepository) FindByID(ctx context.Context, id kernel.UserID) (*domain.UserProfile, error) {
	var profileModel model.UserProfile
	if err := gormutils.TranslateError(repo.db.WithContext(ctx).First(&profileModel, "id = ?", id).Error); err != nil {
		return nil, err
	}
	return repo.toDomain(&profileModel), nil
}

func (repo *UserProfileRepository) toModel(profile *domain.UserProfile) *model.UserProfile {
	return &model.UserProfile{
		ID:          profile.ID(),
		Name:        profile.Name(),
		Gender:      profile.Gender(),
		Email:       profile.Email(),
		PhoneNumber: profile.PhoneNumber(),
		Address:     profile.Address(),
		Sign:        profile.Sign(),
		SignedUpAt:  profile.SignedUpAt(),
	}
}

func (repo *UserProfileRepository) toDomain(profile *model.UserProfile) *domain.UserProfile {
	return domain.LoadUserProfile(
		profile.ID,
		profile.Name,
		profile.Gender,
		profile.Email,
		profile.PhoneNumber,
		profile.Address,
		profile.Sign,
		profile.SignedUpAt,
	)
}
