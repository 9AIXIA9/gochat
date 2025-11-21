package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm/clause"
)

var _ domain.PrivateMessageRepository = (*PrivateMessageRepository)(nil)

type PrivateMessageRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewPrivateMessageRepository(unitOfWork *gormutils.UnitOfWork) *PrivateMessageRepository {
	return &PrivateMessageRepository{unitOfWork: unitOfWork}
}

func (repo *PrivateMessageRepository) SavePrivateMessage(ctx context.Context, message *domain.PrivateMessage) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"sender_id", "recipient_id", "state", "content", "sent_at"}),
		},
	).Create(&model.PrivateMessage{
		ID:          message.ID(),
		SenderID:    message.SenderID(),
		RecipientID: message.RecipientID(),
		State:       message.State(),
		Content:     message.Content(),
		SentAt:      message.SentAt(),
	}).Error)
}

func (repo *PrivateMessageRepository) FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*domain.PrivateMessage, error) {
	var m model.PrivateMessage
	if err := repo.unitOfWork.DB(ctx).WithContext(ctx).Where("id = ?", messageID).First(&m).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.LoadPrivateMessage(
		m.ID,
		m.SenderID,
		m.RecipientID,
		m.State,
		m.Content,
		m.SentAt,
	), nil
}

func (repo *PrivateMessageRepository) FindReceivedPrivateMessage(ctx context.Context, userID kernel.UserID) ([]*domain.PrivateMessage, error) {
	var messages []model.PrivateMessage
	if err := repo.unitOfWork.DB(ctx).WithContext(ctx).Where("recipient_id = ? AND state = ?", userID, domain.MessageStateReceived).Find(&messages).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}

	domainMessages := make([]*domain.PrivateMessage, 0, len(messages))
	for _, message := range messages {
		domainMessages = append(domainMessages, domain.LoadPrivateMessage(
			message.ID,
			message.SenderID,
			message.RecipientID,
			message.State,
			message.Content,
			message.SentAt,
		))
	}
	return domainMessages, nil
}
