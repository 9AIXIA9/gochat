package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/infrastructure/persistence/model"
	"gochat/internal/social/port/kafka"

	"gorm.io/gorm"
)

var _ kafka.UserIDSaver = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) SaveID(ctx context.Context, id kernel.UserID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&model.User{
		ID: id,
	}).Error)
}
