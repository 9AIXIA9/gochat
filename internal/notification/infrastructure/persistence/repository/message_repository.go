package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ application.MessageRepository = (*MessageRepository)(nil)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// SaveMessage 持久化消息及其针对接收者的状态
func (m *MessageRepository) SaveMessage(ctx context.Context, recipient kernel.UserID, message *domain.Message) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 持久化消息主体
		msgModel := &model.Message{
			Model: gorm.Model{
				CreatedAt: message.SentAt(),
			},
			ID:       message.ID(),
			Content:  message.Content(),
			SenderID: message.Sender(),
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(msgModel).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		stateModel := &model.MessageState{
			MessageID: message.ID(),
			UserID:    recipient,
			State:     message.State(),
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(stateModel).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}
