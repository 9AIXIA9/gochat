package repository

import (
	"context"
	"gochat/internal/profile/domain"
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
	//TODO implement me
	panic("implement me")
}

func (repo *UserProfileRepository) Update(ctx context.Context, profile *domain.UserProfile) error {
	//TODO implement me
	panic("implement me")
}

func (repo *UserProfileRepository) FindByID(ctx context.Context, id kernel.UserID) (*domain.UserProfile, error) {
	//TODO implement me
	panic("implement me")
}
