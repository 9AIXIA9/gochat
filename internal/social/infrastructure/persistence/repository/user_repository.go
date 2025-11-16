package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
	"gochat/internal/social/infrastructure/persistence/model"

	"gorm.io/gorm/clause"
)

var _ application.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewUserRepository(unitOfWork *gormutils.UnitOfWork) *UserRepository {
	return &UserRepository{unitOfWork: unitOfWork}
}

func (repo *UserRepository) SaveID(ctx context.Context, id kernel.UserID) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
			DoNothing: true,
		}).Create(&model.User{
		ID: id,
	}).Error)
}
