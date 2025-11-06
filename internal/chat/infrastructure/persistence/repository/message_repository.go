package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ application.PrivateMessageSaver = (*MessageRepository)(nil)
var _ application.RoomMessageSaver = (*MessageRepository)(nil)
var _ application.MessageStateUpdater = (*MessageRepository)(nil)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repo *MessageRepository) SavePrivateMessages(ctx context.Context, recipient kernel.UserID, message []*domain.Message) error {
	//TODO implement me
	panic("implement me")
}

func (repo *MessageRepository) SaveRoomMessages(ctx context.Context, roomID domain.RoomID, message []*domain.Message) error {
	//TODO implement me
	panic("implement me")
}

func (repo *MessageRepository) Update(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.MessageState) error {
	//TODO implement me
	panic("implement me")
}
