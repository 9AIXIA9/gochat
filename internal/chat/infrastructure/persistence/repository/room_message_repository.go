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

		txCtx := context.WithValue(ctx, "transaction", tx)

		if err := repo.eventRepo.CreateUnpublishedEvents(txCtx, message.GetEvents()); err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *RoomMessageRepository) FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*domain.RoomMessage, error) {
	var message model.RoomMessage
	if err := repo.db.WithContext(ctx).
		Where("id = ?", messageID).
		First(&message).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&message), nil
}

func (repo *RoomMessageRepository) toModel(message *domain.RoomMessage) *model.RoomMessage {
	return &model.RoomMessage{
		ID:       message.ID(),
		Content:  message.Content(),
		RoomID:   message.RoomID(),
		SenderID: message.SenderID(),
		SentAt:   message.SentAt(),
	}
}

func (repo *RoomMessageRepository) toDomain(message *model.RoomMessage) *domain.RoomMessage {
	return domain.LoadRoomMessage(message.ID, message.SenderID, message.RoomID, message.Content, message.SentAt)
}
