package repository

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gochat/internal/utils"
	"gorm.io/gorm"
)

var _ domain.MessageRepository = (*MessageRepository)(nil)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (m *MessageRepository) SaveMessage(ctx context.Context, msg *domain.Message) error {
	return utils.HandleDatabaseError(ctx, m.db.WithContext(ctx).Create(model.MessageFromDomain(msg)).Error)
}

func (m *MessageRepository) UpdateMessageSent(ctx context.Context, userNumber domain.UserNumber, msgID domain.MessageID) error {
	return utils.HandleDatabaseError(ctx, m.db.WithContext(ctx).
		Model(&model.UserMessage{}).
		Where(" message_id = ? AND user_number = ?", msgID, userNumber).
		Update("sent", true).Error)
}

func (m *MessageRepository) FindUnsentMessages(ctx context.Context, number domain.UserNumber) ([]*domain.Message, error) {
	var userMsgs []model.UserMessage
	if err := m.db.WithContext(ctx).
		Where("user_number = ? AND sent = ?", number, false).
		Find(&userMsgs).Error; err != nil {
		return nil, utils.HandleDatabaseError(ctx, err)
	}
	if len(userMsgs) == 0 {
		return nil, nil
	}

	// 提取 message_id
	msgIDs := make([]domain.MessageID, 0, len(userMsgs))
	for _, um := range userMsgs {
		msgIDs = append(msgIDs, um.MessageID)
	}

	// 查询消息内容
	var msgs []*model.Message
	if err := m.db.WithContext(ctx).
		Where("id IN ?", msgIDs).
		Order("sent_at ASC").
		Find(&msgs).Error; err != nil {
		return nil, utils.HandleDatabaseError(ctx, err)
	}

	return model.ToDomainMessages(msgs), nil
}
