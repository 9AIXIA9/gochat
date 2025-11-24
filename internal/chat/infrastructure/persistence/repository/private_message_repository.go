package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.PrivateMessageRepository = (*PrivateMessageRepository)(nil)

type PrivateMessageRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewPrivateMessageRepository(db *gorm.DB, eventRepo event.Repository) *PrivateMessageRepository {
	return &PrivateMessageRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *PrivateMessageRepository) Create(ctx context.Context, message *domain.PrivateMessage) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := repo.db.WithContext(ctx).Create(repo.toModel(message)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		txCtx := context.WithValue(ctx, gormutils.UnitOfWorkKey, tx)

		if err := repo.eventRepo.CreateUnpublishedEvents(txCtx, message.GetEvents()); err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *PrivateMessageRepository) FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*domain.PrivateMessage, error) {
	var message model.PrivateMessage
	if err := repo.db.WithContext(ctx).
		Where("id = ?", messageID).
		First(&message).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&message), nil
}

func (repo *PrivateMessageRepository) toModel(message *domain.PrivateMessage) *model.PrivateMessage {
	return &model.PrivateMessage{
		ID:          message.ID(),
		Content:     message.Content(),
		RecipientID: message.RecipientID(),
		SenderID:    message.SenderID(),
		SentAt:      message.SentAt(),
	}
}

func (repo *PrivateMessageRepository) toDomain(message *model.PrivateMessage) *domain.PrivateMessage {
	return domain.LoadPrivateMessage(message.ID, message.SenderID, message.RecipientID, message.Content, message.SentAt)
}
