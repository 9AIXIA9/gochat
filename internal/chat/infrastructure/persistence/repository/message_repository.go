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

type MessageRepository struct {
	innerRepository *gormutils.Repository[model.MessageInformation, domain.MessageInformation]
	db              *gorm.DB
	converter       converter.MessageConverter
}

func NewMessageRepository(db *gorm.DB, converter gormutils.GenericModelConverter[*model.MessageInformation, *domain.MessageInformation]) *MessageRepository {
	return &MessageRepository{innerRepository: gormutils.NewRepository(db, converter), db: db}
}

func (repo *MessageRepository) SavePrivateMessages(ctx context.Context, messages []*domain.PrivateMessage) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, message := range messages {
			if err := repo.innerRepository.WithTx(tx).Save(ctx, message.MessageInformation); err != nil {
				return err
			}

			messageStateModel := repo.converter.PrivateMessageToState(message)
			if err := tx.Create(messageStateModel).Error; err != nil {
				return gormutils.TranslateError(err)
			}
		}
		return nil
	})
}

func (repo *MessageRepository) Update(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.MessageState) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Model(&model.MessageState{}).Where("message_id = ? AND recipient = ?", messageID, recipientID).
		Update("state", newState).Error)
}
