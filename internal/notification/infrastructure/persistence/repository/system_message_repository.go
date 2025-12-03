package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.SystemMessageRepository = (*SystemMessageRepository)(nil)

type SystemMessageRepository struct {
	db *gorm.DB
}

func NewSystemMessageRepository(db *gorm.DB) *SystemMessageRepository {
	return &SystemMessageRepository{db: db}
}

func (repo *SystemMessageRepository) Create(ctx context.Context, message *domain.SystemMessage) error {
	if err := repo.db.WithContext(ctx).Create(repo.toModel(message)).Error; err != nil {
		return gormutils.TranslateError(err)
	}
	return nil
}

func (repo *SystemMessageRepository) Update(ctx context.Context, message *domain.SystemMessage) error {
	if err := repo.db.WithContext(ctx).Updates(repo.toModel(message)).Error; err != nil {
		return gormutils.TranslateError(err)
	}
	return nil
}

func (repo *SystemMessageRepository) Updates(ctx context.Context, messages []*domain.SystemMessage) error {
	for _, message := range messages {
		if err := repo.db.WithContext(ctx).Updates(repo.toModel(message)).Error; err != nil {
			return gormutils.TranslateError(err)
		}
	}
	return nil
}

func (repo *SystemMessageRepository) FindByID(ctx context.Context, messageID kernel.MessageID) (*domain.SystemMessage, error) {
	var message model.SystemMessage
	if err := repo.db.WithContext(ctx).First(&message, "id = ?", messageID).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&message), nil
}

func (repo *SystemMessageRepository) FindsByState(ctx context.Context, userID kernel.UserID, state domain.MessageState) ([]*domain.SystemMessage, error) {
	var messages []model.SystemMessage
	if err := repo.db.WithContext(ctx).Find(&messages, "recipient_id = ? AND state = ?", userID, state).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomains(messages), nil
}

func (repo *SystemMessageRepository) toModel(message *domain.SystemMessage) *model.SystemMessage {
	return &model.SystemMessage{
		ID:          message.ID(),
		Topic:       message.Topic(),
		RecipientID: message.RecipientID(),
		State:       message.State(),
		Content:     message.Content(),
		SentAt:      message.SentAt(),
	}
}

func (repo *SystemMessageRepository) toModels(messages []*domain.SystemMessage) []*model.SystemMessage {
	models := make([]*model.SystemMessage, len(messages))
	for i, msg := range messages {
		models[i] = repo.toModel(msg)
	}
	return models
}

func (repo *SystemMessageRepository) toDomain(message *model.SystemMessage) *domain.SystemMessage {
	return domain.LoadSystemMessage(message.ID, message.Topic, message.RecipientID, message.State, message.Content, message.SentAt)
}

func (repo *SystemMessageRepository) toDomains(messages []model.SystemMessage) []*domain.SystemMessage {
	domains := make([]*domain.SystemMessage, len(messages))
	for i, msg := range messages {
		domains[i] = repo.toDomain(&msg)
	}
	return domains
}
