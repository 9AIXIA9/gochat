package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
)

var _ domain.PrivateMessageRepository = (*PrivateMessageRepository)(nil)

type PrivateMessageRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewPrivateMessageRepository(unitOfWork *gormutils.UnitOfWork) *PrivateMessageRepository {
	return &PrivateMessageRepository{unitOfWork: unitOfWork}
}

func (repo *PrivateMessageRepository) SavePrivateMessage(ctx context.Context, message *domain.PrivateMessage) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&model.PrivateMessage{
		ID:          message.ID(),
		Content:     message.Content(),
		RecipientID: message.RecipientID(),
		SenderID:    message.SenderID(),
		SentAt:      message.SentAt(),
	}).Error)
}

func (repo *PrivateMessageRepository) FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*domain.PrivateMessage, error) {
	var m model.PrivateMessage
	err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Where("id = ?", messageID).
		First(&m).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.LoadPrivateMessage(m.ID, m.SenderID, m.RecipientID, m.Content, m.SentAt), nil
}
