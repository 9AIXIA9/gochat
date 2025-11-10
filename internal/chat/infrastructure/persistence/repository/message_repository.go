package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ application.MessageRepository = (*MessageRepository)(nil)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repo *MessageRepository) SavePrivateMessage(ctx context.Context, recipient kernel.UserID, message *domain.Message) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&model.PrivateMessage{
		Model: gorm.Model{
			CreatedAt: message.SentAt(),
		},
		ID:          message.ID(),
		Content:     message.Content(),
		RecipientID: recipient,
		SenderID:    message.Sender(),
	}).Error)
}

func (repo *MessageRepository) SaveRoomMessage(ctx context.Context, roomID domain.RoomID, message *domain.Message) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(&model.RoomMessage{
		Model:    gorm.Model{},
		ID:       message.ID(),
		Content:  message.Content(),
		RoomID:   roomID,
		SenderID: message.Sender(),
	}).Error)
}
