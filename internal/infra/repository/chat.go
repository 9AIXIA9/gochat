package repository

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
	"time"
)

type ChatRepository struct {
	db *gorm.DB
}

type Chat struct {
	gorm.Model
	UserNumber domain.UserNumber
	RoomNumber domain.RoomNumber
	Content    string
	SendTime   time.Time
}

func NewChatRepository(db *gorm.DB) domain.ChatRepository {
	return &ChatRepository{db: db}
}
