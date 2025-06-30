package repository

import (
	"context"
	"errors"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, user *domain.User) (bool, error) {
	var existingUser model.User
	result := u.db.WithContext(ctx).Where("number = ?", user.Number).First(&existingUser)
	if result.RowsAffected > 0 {
		return true, nil // 用户已存在
	}

	gormUser := model.UserFromDomain(user)
	if err := u.db.WithContext(ctx).Create(gormUser).Error; err != nil {
		return false, err
	}

	return false, nil
}

func (u *UserRepository) QueryByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error) {
	var gormUser model.User
	result := u.db.WithContext(ctx).Where("number = ?", number).First(&gormUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, result.Error
	}

	return gormUser.ToDomain(), nil
}
