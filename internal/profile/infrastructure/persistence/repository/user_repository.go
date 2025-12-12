package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/profile/domain"
	"gochat/internal/profile/infrastructure/persistence/model"

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
	return gormutils.TranslateError(repo.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				DoNothing: true,
			},
		).
		Create(repo.toModel(user)).Error)
}

func (repo *UserRepository) toModel(user *domain.User) *model.User {
	return &model.User{
		ID: user.ID(),
	}
}
