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

var _ domain.RoomMessageRepository = (*RoomMessageRepository)(nil)

type RoomMessageRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewRoomMessageRepository(db *gorm.DB, eventRepo event.Repository) *RoomMessageRepository {
	return &RoomMessageRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *RoomMessageRepository) Create(ctx context.Context, message *domain.RoomMessage) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(repo.toModel(message)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		txCtx := gormutils.SetTransaction(ctx, tx)

		if err := repo.eventRepo.CreateUnpublishedEvents(txCtx, message.GetEvents()); err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *RoomMessageRepository) FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*domain.RoomMessage, error) {
	var message model.RoomMessage
	if err := repo.db.WithContext(ctx).
		Preload("Recipients").
		Where("id = ?", messageID).
		First(&message).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&message), nil
}

func (repo *RoomMessageRepository) FindsByRoomIDAndUserIDWithoutRecipients(ctx context.Context, roomID kernel.RoomID, userID kernel.UserID, limit int, baseID kernel.MessageID) ([]*domain.RoomMessage, error) {
	var messages []model.RoomMessage
	query := repo.db.WithContext(ctx).
		Model(&model.RoomMessage{}).
		Joins("JOIN chat_roomships ON chat_roomships.room_id = chat_room_messages.room_id AND chat_roomships.user_id = ?", userID).
		Where("room_id = ?", roomID).
		Order("id DESC").
		Limit(limit)

	if baseID != "" {
		query = query.Where("id < ?", baseID)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomainsWithoutRecipients(messages), nil
}

func (repo *RoomMessageRepository) toModel(message *domain.RoomMessage) *model.RoomMessage {
	recipients := make([]*model.RoomMessageRecipient, len(message.RecipientIDs()))
	for i, recipientID := range message.RecipientIDs() {
		recipients[i] = &model.RoomMessageRecipient{
			MessageID: message.ID(),
			UserID:    recipientID,
		}
	}

	return &model.RoomMessage{
		ID:         message.ID(),
		SenderID:   message.SenderID(),
		RoomID:     message.RoomID(),
		Content:    message.Content(),
		SentAt:     message.SentAt(),
		Recipients: recipients,
	}
}

func (repo *RoomMessageRepository) toDomain(message *model.RoomMessage) *domain.RoomMessage {
	recipients := make([]kernel.UserID, len(message.Recipients))
	for i, recipient := range message.Recipients {
		recipients[i] = recipient.UserID
	}

	return domain.LoadRoomMessage(
		message.ID,
		message.SenderID,
		recipients,
		message.RoomID,
		message.Content,
		message.SentAt,
	)
}

func (repo *RoomMessageRepository) toDomains(messages []model.RoomMessage) []*domain.RoomMessage {
	domainMessages := make([]*domain.RoomMessage, 0, len(messages))
	for _, message := range messages {
		domainMessages = append(domainMessages, repo.toDomain(&message))
	}
	return domainMessages
}

func (repo *RoomMessageRepository) toDomainWithoutRecipients(message *model.RoomMessage) *domain.RoomMessage {
	return domain.LoadRoomMessage(
		message.ID,
		message.SenderID,
		nil,
		message.RoomID,
		message.Content,
		message.SentAt,
	)
}

func (repo *RoomMessageRepository) toDomainsWithoutRecipients(messages []model.RoomMessage) []*domain.RoomMessage {
	domainMessages := make([]*domain.RoomMessage, 0, len(messages))
	for _, message := range messages {
		domainMessages = append(domainMessages, repo.toDomainWithoutRecipients(&message))
	}
	return domainMessages
}
