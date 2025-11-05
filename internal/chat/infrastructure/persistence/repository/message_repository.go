package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/converter"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ application.PrivateMessagesSaver = (*MessageRepository)(nil)
var _ application.MessageStateUpdater = (*MessageRepository)(nil)
var _ application.RoomMessagesSaver = (*MessageRepository)(nil)

type MessageRepository struct {
	db        *gorm.DB
	converter *converter.MessageConverter
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		db:        db,
		converter: &converter.MessageConverter{},
	}
}

func (repo *MessageRepository) SavePrivateMessages(ctx context.Context, messages []*domain.PrivateMessage) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, message := range messages {
			info := repo.converter.ToInformation(message.ID(), message.MessageInformation)
			if err := tx.Create(info).Error; err != nil {
				return err
			}

			state := repo.converter.ToRecipientState(message.ID(), message.RecipientMessageState)
			if err := tx.Create(state).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *MessageRepository) SaveRoomMessages(ctx context.Context, messages []*domain.RoomMessage) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, message := range messages {
			info := repo.converter.ToInformation(message.ID(), message.MessageInformation)
			if err := tx.Create(info).Error; err != nil {
				return err
			}

			state := repo.converter.ToRecipientStates(message.ID(), message.States())
			if err := tx.Create(state).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *MessageRepository) Update(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.State) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).
		Model(&model.RecipientMessageState{}).
		Where("message_id = ? AND recipient = ?", messageID, recipientID).
		Update("state", newState).Error)
}
