package repository

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gochat/internal/types"
	"gochat/internal/utils"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) domain.ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Save(ctx context.Context, chat *domain.Chat) error {
	gormChat := model.ChatFromDomain(chat)
	if err := r.db.WithContext(ctx).Save(gormChat).Error; err != nil {
		return utils.CheckDuplicateKeyError(err)
	}
	return nil
}

func (r *ChatRepository) FindAllByRoomNumber(ctx context.Context, number domain.RoomNumber) ([]*domain.Chat, error) {
	var gormChats []model.Chat
	if err := r.db.WithContext(ctx).Where("room_number = ?", number).Order("sent_at").Find(&gormChats).Error; err != nil {
		return nil, utils.CheckNotFoundError(err)
	}

	if len(gormChats) == 0 {
		return nil, types.ErrNotFound
	}

	chats := make([]*domain.Chat, len(gormChats))
	for i, gormChat := range gormChats {
		chats[i] = gormChat.ToDomain()
	}

	return chats, nil
}
