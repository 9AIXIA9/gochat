package repository

import (
	"context"
	"gochat/internal/profile/domain"

	"gorm.io/gorm"
)

var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Save(ctx context.Context, user *domain.User) error {
	//TODO implement me
	panic("implement me")
}
