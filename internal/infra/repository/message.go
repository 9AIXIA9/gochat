package repository

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/model"
	"gochat/internal/utils"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &MessageRepository{db: db}
}

func (m *MessageRepository) Save(ctx context.Context, message *domain.Message) error {
	return utils.HandleDatabaseError(ctx, m.db.WithContext(ctx).Create(model.MessageFromDomain(message)).Error)
}

func (m *MessageRepository) UpdateMessageSent(ctx context.Context, userNumber domain.UserNumber, msgID domain.MessageID) error {
	return utils.HandleDatabaseError(ctx, m.db.WithContext(ctx).
		Model(&model.UserMessage{}).
		Where(" message_id = ? AND user_number = ?", msgID, userNumber).
		Update("sent", true).Error)
}

func (m *MessageRepository) SaveAndQueryUserNumberShouldSent(ctx context.Context, message *domain.Message) ([]domain.UserNumber, error) {
	var userNumbers []domain.UserNumber
	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 存储消息本体
		modelMsg := model.MessageFromDomain(message)
		if err := tx.Create(modelMsg).Error; err != nil {
			return utils.HandleDatabaseError(ctx, err)
		}

		// 查询房间成员
		var userRooms []model.UserRoom
		if err := tx.Where("room_number = ?", message.To()).Find(&userRooms).Error; err != nil {
			return utils.HandleDatabaseError(ctx, err)
		}

		if len(userRooms) > 0 {
			for _, ur := range userRooms {
				userNumbers = append(userNumbers, ur.UserNumber)
			}
		} else {
			// 校验用户号是否存在
			var user model.User
			if err := tx.Where("number = ?", message.To()).First(&user).Error; err != nil {
				return utils.HandleDatabaseError(ctx, err) // 用户不存在或查询出错
			}
			userNumbers = append(userNumbers, domain.UserNumber(message.To()))
		}

		// 批量插入 UserMessage
		userMessages := make([]*model.UserMessage, 0, len(userNumbers))
		for _, num := range userNumbers {
			userMessages = append(userMessages, &model.UserMessage{
				Sent:       false,
				UserNumber: num,
				MessageID:  message.ID(),
			})
		}
		if len(userMessages) > 0 {
			if err := tx.Create(&userMessages).Error; err != nil {
				return utils.HandleDatabaseError(ctx, err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, utils.HandleDatabaseError(ctx, err)
	}
	return userNumbers, nil
}

func (m *MessageRepository) UpdateMessagesSentToOneUser(ctx context.Context, number domain.UserNumber, msgIDs []domain.MessageID) error {
	if len(msgIDs) == 0 {
		return nil
	}
	if err := m.db.WithContext(ctx).
		Model(&model.UserMessage{}).
		Where("user_number = ? AND message_id IN ?", number, msgIDs).
		Update("sent", true).Error; err != nil {
		return utils.HandleDatabaseError(ctx, err)
	}
	return nil
}

func (m *MessageRepository) UpdateMessageSentToManyUsers(ctx context.Context, msgID domain.MessageID, userNumbers []domain.UserNumber) error {
	if len(userNumbers) == 0 {
		return nil
	}

	if err := m.db.WithContext(ctx).
		Model(&model.UserMessage{}).
		Where(" message_id = ? AND user_number in ?", msgID, userNumbers).
		Update("sent", true).Error; err != nil {
		return utils.HandleDatabaseError(ctx, err)
	}
	return nil
}

func (m *MessageRepository) QueryUnsentMessages(ctx context.Context, number domain.UserNumber) ([]*domain.Message, error) {
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
