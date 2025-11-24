package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/social/domain"
	"gochat/internal/social/infrastructure/persistence/model"

	"gorm.io/gorm/clause"
)

var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewUserRepository(unitOfWork *gormutils.UnitOfWork) *UserRepository {
	return &UserRepository{unitOfWork: unitOfWork}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
			DoNothing: true,
		}).Create(repo.toModel(user)).Error)
}

func (repo *UserRepository) toModel(user *domain.User) *model.User {
	return &model.User{
		ID: user.ID(),
	}
}
