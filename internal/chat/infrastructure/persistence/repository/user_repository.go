package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
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

func (repo *UserRepository) FindByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error) {
	var user model.User
	err := repo.db.WithContext(ctx).First(&user, "number = ?", number).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}

	return domain.NewUser(user.ID, user.Number), nil
}

func (repo *UserRepository) SaveNumber(ctx context.Context, userID kernel.UserID, number domain.UserNumber) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&model.User{
		ID:     userID,
		Number: number,
	}).Error)
}
