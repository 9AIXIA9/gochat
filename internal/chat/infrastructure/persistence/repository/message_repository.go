package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ application.MessageSaver = (*MessageRepository)(nil)
var _ application.MessageStateUpdater = (*MessageRepository)(nil)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (m *MessageRepository) Saves(ctx context.Context, message []*domain.Message) error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageRepository) Update(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.MessageState) error {
	//TODO implement me
	panic("implement me")
}
