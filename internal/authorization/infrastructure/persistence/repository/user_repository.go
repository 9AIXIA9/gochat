package repository

import (
	"context"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"

	"gorm.io/gorm"
)

var _ application.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	innerRepository *gormutils.Repository[model.User, domain.User]
	db              *gorm.DB
}

func NewUserRepository(db *gorm.DB, converter gormutils.GenericModelConverter[*model.User, *domain.User]) *UserRepository {
	return &UserRepository{innerRepository: gormutils.NewRepository(db, converter), db: db}
}

func (repo *UserRepository) Save(ctx context.Context, user *domain.User) error {
	return repo.innerRepository.Save(ctx, user)
}

func (repo *UserRepository) FindByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error) {
	return repo.innerRepository.Find(ctx, gormutils.Where("number = ?", number))
}
