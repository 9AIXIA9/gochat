package repository

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gochat/internal/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Save(ctx context.Context, user *domain.User) error {
	gormUser := model.UserFromDomain(user)
	if err := u.db.WithContext(ctx).Save(gormUser).Error; err != nil {
		return utils.HandleDatabaseError(ctx, err)
	}
	return nil
}

func (u *UserRepository) FindOneByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error) {
	var gormUser model.User
	if err := u.db.WithContext(ctx).Where("number = ?", number).First(&gormUser).Error; err != nil {
		return nil, utils.HandleDatabaseError(ctx, err)
	}

	return gormUser.ToDomain(), nil
}
