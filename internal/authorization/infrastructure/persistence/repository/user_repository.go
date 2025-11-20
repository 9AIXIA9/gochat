package repository

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
)

var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewUserRepository(unitOfWork *gormutils.UnitOfWork) *UserRepository {
	return &UserRepository{unitOfWork: unitOfWork}
}

func (repo *UserRepository) Save(ctx context.Context, user *domain.User) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&model.User{
		ID:                user.ID(),
		Email:             user.Email(),
		Number:            user.Number(),
		PasswordEncrypted: user.PasswordEncrypted(),
		SignedUpAt:        user.SignedUpAt(),
	}).Error)
}

func (repo *UserRepository) FindByNumber(ctx context.Context, number kernel.UserNumber) (*domain.User, error) {
	var user model.User
	if err := repo.unitOfWork.DB(ctx).WithContext(ctx).First(&user, "number = ?", number).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.LoadUser(user.ID, user.Email, user.Number, user.PasswordEncrypted, user.SignedUpAt), nil
}

func (repo *UserRepository) FindByID(ctx context.Context, id kernel.UserID) (*domain.User, error) {
	var user model.User
	if err := repo.unitOfWork.DB(ctx).WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.LoadUser(user.ID, user.Email, user.Number, user.PasswordEncrypted, user.SignedUpAt), nil
}
