package repository

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Number  domain.UserNumber
	Name    string
	PwdHash string
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(user *domain.User) (bool, error) {
	return false, nil
}

func (u *UserRepository) QueryByNumber(number domain.UserNumber) (*domain.User, error) {
	return nil, nil
}
