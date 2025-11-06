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

var _ application.MessageSaver = (*MessageRepository)(nil)
var _ application.MessageStateUpdater = (*MessageRepository)(nil)

type MessageRepository struct {
	db              *gorm.DB
	innerRepository *gormutils.Repository[model.Message, domain.Message]
}

func NewMessageRepository(db *gorm.DB, converter gormutils.GenericModelConverter[*model.Message, *domain.Message]) *MessageRepository {
	return &MessageRepository{innerRepository: gormutils.NewRepository(db, converter), db: db}
}

func (repo *MessageRepository) Saves(ctx context.Context, messages []*domain.Message) error {
	return repo.innerRepository.Saves(ctx, messages)
}

func (repo *MessageRepository) Update(ctx context.Context, messageID domain.MessageID, recipientID kernel.UserID, newState domain.MessageState) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).
		Model(&model.RecipientMessageState{}).
		Where("message_id = ? AND recipient = ?", messageID, recipientID).
		Update("state", newState).Error)
}
