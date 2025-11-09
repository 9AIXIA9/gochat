package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/application"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ application.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) SaveEmail(ctx context.Context, userID kernel.UserID, email kernel.Email) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&model.User{
		ID:    userID,
		Email: email,
	}).Error)
}
