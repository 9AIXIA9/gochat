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
		if err := tx.Create(repo.toModel(message)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		txCtx := gormutils.SetTransaction(ctx, tx)

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

func (repo *PrivateMessageRepository) FindsByUserIDs(ctx context.Context, userID1, userID2 kernel.UserID, limit int, baseID kernel.MessageID) ([]*domain.PrivateMessage, error) {
	var messages []model.PrivateMessage

	query := repo.db.WithContext(ctx).Model(&model.PrivateMessage{}).
		Where("(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
			userID1, userID2, userID2, userID1).
		Order("id DESC").
		Limit(limit)
	if baseID != "" {
		query = query.Where("id < ?", baseID)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomains(messages), nil
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
	return domain.LoadPrivateMessage(
		message.ID,
		message.SenderID,
		message.RecipientID,
		message.Content,
		message.SentAt,
	)
}

func (repo *PrivateMessageRepository) toDomains(messages []model.PrivateMessage) []*domain.PrivateMessage {
	domainMessages := make([]*domain.PrivateMessage, 0, len(messages))
	for _, message := range messages {
		domainMessages = append(domainMessages, repo.toDomain(&message))
	}
	return domainMessages
}
