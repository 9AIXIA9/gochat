package repository

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &MessageRepository{db: db}
}
