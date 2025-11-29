package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

//TODO 查询功能仍有 bug 待完善

const maxLimit = 10

var _ domain.PrivateMessageRepository = (*PrivateMessageRepository)(nil)

type PrivateMessageRepository struct {
	db *gorm.DB
}

func NewPrivateMessageRepository(db *gorm.DB) *PrivateMessageRepository {
	return &PrivateMessageRepository{db: db}
}

func (repo *PrivateMessageRepository) Create(ctx context.Context, message *domain.PrivateMessage) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(repo.toModel(message)).Error)
}

func (repo *PrivateMessageRepository) Update(ctx context.Context, message *domain.PrivateMessage) error {
	modelMessage := repo.toModel(message)
	return gormutils.TranslateError(repo.db.WithContext(ctx).Model(&modelMessage).Updates(modelMessage).Error)
}

func (repo *PrivateMessageRepository) Updates(ctx context.Context, messages []*domain.PrivateMessage) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, message := range messages {
			modelMessage := repo.toModel(message)
			if err := tx.Model(&modelMessage).Updates(modelMessage).Error; err != nil {
				return gormutils.TranslateError(err)
			}
		}
		return nil
	})
}

func (repo *PrivateMessageRepository) FindByID(ctx context.Context, messageID kernel.MessageID) (*domain.PrivateMessage, error) {
	var modelMessage model.PrivateMessage
	if err := repo.db.WithContext(ctx).First(&modelMessage, "id = ?", messageID).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&modelMessage), nil
}

func (repo *PrivateMessageRepository) FindsByState(ctx context.Context, userID kernel.UserID, state domain.MessageState) ([]*domain.PrivateMessage, error) {
	var modelMessages []*model.PrivateMessage
	if err := repo.db.WithContext(ctx).
		Where("recipient_id = ? AND state = ?", userID, state).
		Find(&modelMessages).
		Limit(maxLimit).
		Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomains(modelMessages), nil
}

func (repo *PrivateMessageRepository) toModel(message *domain.PrivateMessage) *model.PrivateMessage {
	return &model.PrivateMessage{
		ID:          message.ID(),
		SenderID:    message.SenderID(),
		RecipientID: message.RecipientID(),
		State:       message.State(),
		Content:     message.Content(),
		SentAt:      message.SentAt(),
	}
}

func (repo *PrivateMessageRepository) toModels(messages []*domain.PrivateMessage) []*model.PrivateMessage {
	modelMessages := make([]*model.PrivateMessage, 0, len(messages))
	for _, message := range messages {
		modelMessages = append(modelMessages, repo.toModel(message))
	}
	return modelMessages
}

func (repo *PrivateMessageRepository) toDomain(message *model.PrivateMessage) *domain.PrivateMessage {
	return domain.LoadPrivateMessage(
		message.ID,
		message.SenderID,
		message.RecipientID,
		message.State,
		message.Content,
		message.SentAt,
	)
}

func (repo *PrivateMessageRepository) toDomains(messages []*model.PrivateMessage) []*domain.PrivateMessage {
	domainMessages := make([]*domain.PrivateMessage, 0, len(messages))
	for _, message := range messages {
		domainMessages = append(domainMessages, repo.toDomain(message))
	}
	return domainMessages
}
