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

var _ application.PrivateMessageSaver = (*MessageRepository)(nil)
var _ application.RoomMessageSaver = (*MessageRepository)(nil)
var _ application.MessageStateUpdater = (*MessageRepository)(nil)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repo *MessageRepository) SavePrivateMessages(ctx context.Context, recipient kernel.UserID, messages []*domain.Message) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, message := range messages {
			if err := tx.Create(&model.Message{
				Model: gorm.Model{
					CreatedAt: message.SentAt(),
				},
				ID:       message.ID(),
				Content:  message.Content(),
				SenderID: message.Sender(),
			}).Error; err != nil {
				return gormutils.TranslateError(err)
			}

			if err := tx.Create(&model.MessageUser{
				MessageID: message.ID(),
				UserID:    recipient,
				State:     message.State(),
			}).Error; err != nil {
				return gormutils.TranslateError(err)
			}
		}
		return nil
	})
}

func (repo *MessageRepository) SaveRoomMessages(ctx context.Context, roomID domain.RoomID, messages []*domain.Message) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var room model.Room
		if err := tx.Preload("Members").First(&room, "id = ?", roomID).Error; err != nil {
			return gormutils.TranslateError(err)
		}
		members := room.Members

		txRepo := NewMessageRepository(tx)

		for _, member := range members {
			if err := txRepo.SavePrivateMessages(ctx, member.ID, messages); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *MessageRepository) Update(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.MessageState) error {
	return gormutils.TranslateError(
		repo.db.WithContext(ctx).
			Model(&model.MessageUser{}).
			Where("message_id = ? AND user_id = ?", messageID, recipientID).
			Update("state", newState).Error,
	)
}
