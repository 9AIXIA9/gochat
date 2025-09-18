package repository

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &MessageRepository{db: db}
}

func (m *MessageRepository) Save(ctx context.Context, message *domain.Message) error {
	modelMsg := model.MessageFromDomain(message)
	return m.db.WithContext(ctx).Create(modelMsg).Error
}

func (m *MessageRepository) UpdateMessagesSent(ctx context.Context, messages []*domain.Message) error {
	for _, msg := range messages {
		if err := m.db.WithContext(ctx).
			Model(&model.Message{}).
			Where("id = ?", msg.ID()).
			Update("sent", msg.IsSent()).Error; err != nil {
			return err
		}
	}
	return nil
}

func (m *MessageRepository) QueryMessages(ctx context.Context, number domain.BaseNumber, count int) ([]*domain.Message, error) {
	var modelMsgs []*model.Message
	err := m.db.WithContext(ctx).
		Where("recipient = ?", number).
		Order("sent_at desc").
		Limit(count).
		Find(&modelMsgs).Error
	if err != nil {
		return nil, err
	}
	return model.ToDomainMessages(modelMsgs), nil
}

func (m *MessageRepository) QueryUnsentMessages(ctx context.Context, number domain.UserNumber) ([]*domain.Message, error) {
	// 查询用户加入的所有房间号
	var roomNumbers []domain.RoomNumber
	err := m.db.WithContext(ctx).
		Model(&model.UserRoom{}).
		Where("user_number = ?", number).
		Pluck("room_number", &roomNumbers).Error
	if err != nil {
		return nil, err
	}

	//构建接收者列表（用户号 + 房间号）
	recipients := []domain.BaseNumber{domain.BaseNumber(number)}
	for _, rn := range roomNumbers {
		recipients = append(recipients, domain.BaseNumber(rn))
	}

	// 查询消息（包含点对点消息和群消息）
	var modelMsgs []*model.Message
	err = m.db.WithContext(ctx).
		Where("recipient IN (?) AND sent = ?", recipients, false).
		Order("sent_at desc").
		Find(&modelMsgs).Error
	if err != nil {
		return nil, err
	}

	return model.ToDomainMessages(modelMsgs), nil
}

func (m *MessageRepository) QueryAllMessages(ctx context.Context, number domain.BaseNumber) ([]*domain.Message, error) {
	var modelMsgs []*model.Message
	err := m.db.WithContext(ctx).
		Where("recipient = ?", number).
		Order("sent_at desc").
		Find(&modelMsgs).Error
	if err != nil {
		return nil, err
	}
	return model.ToDomainMessages(modelMsgs), nil
}
