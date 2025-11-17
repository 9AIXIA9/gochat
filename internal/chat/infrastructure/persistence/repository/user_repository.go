package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm/clause"
)

var _ application.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewUserRepository(unitOfWork *gormutils.UnitOfWork) *UserRepository {
	return &UserRepository{unitOfWork: unitOfWork}
}

func (repo *UserRepository) FindByNumber(ctx context.Context, number kernel.UserNumber) (*domain.User, error) {
	var user model.User
	err := repo.unitOfWork.DB(ctx).WithContext(ctx).First(&user, "number = ?", number).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}

	return domain.NewUser(user.ID, user.Number), nil
}

func (repo *UserRepository) SaveNumber(ctx context.Context, userID kernel.UserID, number kernel.UserNumber) error {
	return gormutils.TranslateError(
		repo.unitOfWork.DB(ctx).WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
			DoNothing: true,
		}).Create(&model.User{
			ID:     userID,
			Number: number,
		}).Error,
	)
}
