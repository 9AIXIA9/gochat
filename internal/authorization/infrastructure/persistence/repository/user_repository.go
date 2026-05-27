package repository

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewUserRepository(db *gorm.DB, eventRepo event.Repository) *UserRepository {
	return &UserRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(repo.toModel(user)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		txCtx := gormutils.SetTransaction(ctx, tx)

		if err := repo.eventRepo.CreateUnpublishedEvents(txCtx, user.GetEvents()); err != nil {
			return gormutils.TranslateError(err)
		}
		return nil
	})
}

func (repo *UserRepository) FindByNumber(ctx context.Context, number kernel.UserNumber) (*domain.User, error) {
	var user model.User
	if err := repo.db.WithContext(ctx).First(&user, "number = ?", number).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&user), nil
}

func (repo *UserRepository) FindByID(ctx context.Context, id kernel.UserID) (*domain.User, error) {
	var user model.User
	if err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&user), nil
}

func (repo *UserRepository) toModel(user *domain.User) *model.User {
	return &model.User{
		ID:                user.ID(),
		Email:             user.Email(),
		Number:            user.Number(),
		PasswordEncrypted: user.PasswordEncrypted(),
		SignedUpAt:        user.SignedUpAt(),
	}
}

func (repo *UserRepository) toDomain(user *model.User) *domain.User {
	return domain.LoadUser(user.ID, user.Email, user.Number, user.PasswordEncrypted, user.SignedUpAt)
}
