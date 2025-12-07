package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/infrastructure/persistence/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Save(ctx context.Context, user *domain.User) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Clauses(
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
