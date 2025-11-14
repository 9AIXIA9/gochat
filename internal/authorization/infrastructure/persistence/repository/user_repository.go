package repository

import (
	"context"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
	"time"

	"gorm.io/gorm"
)

var _ application.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Save(ctx context.Context, user *domain.User) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&model.User{
		ID:                user.ID(),
		Email:             user.Email(),
		Number:            user.Number(),
		PasswordEncrypted: user.PasswordEncrypted(),
		LastLoggedInAt:    user.LastLoggedInAt(),
		CreatedAt:         user.SignedUpAt(),
	}).Error)
}

func (repo *UserRepository) FindByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error) {
	var user model.User
	if err := repo.db.WithContext(ctx).First(&user, "number = ?", number).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.NewUser(user.ID, user.Email, user.Number, user.PasswordEncrypted, user.LastLoggedInAt, user.CreatedAt), nil
}

func (repo *UserRepository) UpdateLoggedInAt(ctx context.Context, userID kernel.UserID, t time.Time) error {
	return gormutils.TranslateError(
		repo.db.WithContext(ctx).
			Model(&model.User{}).
			Where("id = ?", userID).
			UpdateColumn("last_logged_in_at", t).Error,
	)
}
func (repo *UserRepository) FindByID(ctx context.Context, id kernel.UserID) (*domain.User, error) {
	var user model.User
	if err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.NewUser(user.ID, user.Email, user.Number, user.PasswordEncrypted, user.LastLoggedInAt, user.CreatedAt), nil
}
