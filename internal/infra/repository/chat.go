package repository

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) domain.ChatRepository {
	return &ChatRepository{db: db}
}

func (c *ChatRepository) Create(ctx context.Context, chat domain.Chat) (bool, error) {
	gormChat := model.ChatFromDomain(chat)
	if err := c.db.WithContext(ctx).Create(gormChat).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (c *ChatRepository) QueryAllByRoomNumber(ctx context.Context, number domain.RoomNumber) ([]domain.Chat, error) {
	var gormChats []model.Chat
	if err := c.db.WithContext(ctx).Where("room_number = ?", number).Order("sent_at").Find(&gormChats).Error; err != nil {
		return nil, err
	}

	chats := make([]domain.Chat, len(gormChats))
	for i, gormChat := range gormChats {
		chats[i] = gormChat.ToDomain()
	}

	return chats, nil
}
