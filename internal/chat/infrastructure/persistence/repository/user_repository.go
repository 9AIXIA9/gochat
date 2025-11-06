package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"

	"gorm.io/gorm"
)

var _ application.UserFinder = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) FindByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}
