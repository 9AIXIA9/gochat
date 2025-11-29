package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomMessageRepository = (*RoomMessageRepository)(nil)

type RoomMessageRepository struct {
	db *gorm.DB
}

func NewRoomMessageRepository(db *gorm.DB) *RoomMessageRepository {
	return &RoomMessageRepository{db: db}
}

func (repo *RoomMessageRepository) Create(ctx context.Context, message *domain.RoomMessage) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Create(repo.toModel(message)).Error)
}

func (repo *RoomMessageRepository) Update(ctx context.Context, message *domain.RoomMessage) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		modelMessage := repo.toModel(message)

		// 先删除旧的关联，避免 GORM 将外键置为 NULL
		if err := tx.Where("message_id = ?", modelMessage.ID).Delete(&model.RoomMessageState{}).Error; err != nil {
			return gormutils.TranslateError(err)
		}
		if err := tx.Where("message_id = ?", modelMessage.ID).Delete(&model.RoomMessageRecipient{}).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		// 再插入新的关联
		if len(modelMessage.States) > 0 {
			if err := tx.Create(&modelMessage.States).Error; err != nil {
				return gormutils.TranslateError(err)
			}
		}
		if len(modelMessage.Recipients) > 0 {
			if err := tx.Create(&modelMessage.Recipients).Error; err != nil {
				return gormutils.TranslateError(err)
			}
		}

		// 更新主记录字段（不触发全量 save，避免误操作关联）
		if err := tx.Model(&model.RoomMessage{}).Where("id = ?", modelMessage.ID).Updates(map[string]interface{}{
			"sender_id": modelMessage.SenderID,
			"room_id":   modelMessage.RoomID,
			"content":   modelMessage.Content,
			"sent_at":   modelMessage.SentAt,
		}).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *RoomMessageRepository) Updates(ctx context.Context, messages []*domain.RoomMessage) error {
	if len(messages) == 0 {
		return nil
	}

	modelMessages := repo.toModels(messages)
	messageIDs := make([]kernel.MessageID, 0, len(modelMessages))
	allStates := make([]*model.RoomMessageState, 0)
	allRecipients := make([]*model.RoomMessageRecipient, 0)

	for _, msg := range modelMessages {
		messageIDs = append(messageIDs, msg.ID)
		allStates = append(allStates, msg.States...)
		allRecipients = append(allRecipients, msg.Recipients...)
	}

	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 批量删除旧的关联
		if err := tx.Where("message_id IN ?", messageIDs).Delete(&model.RoomMessageState{}).Error; err != nil {
			return gormutils.TranslateError(err)
		}
		if err := tx.Where("message_id IN ?", messageIDs).Delete(&model.RoomMessageRecipient{}).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		// 批量插入新的关联
		if len(allStates) > 0 {
			if err := tx.Create(&allStates).Error; err != nil {
				return gormutils.TranslateError(err)
			}
		}
		if len(allRecipients) > 0 {
			if err := tx.Create(&allRecipients).Error; err != nil {
				return gormutils.TranslateError(err)
			}
		}

		// 批量更新 RoomMessage 表本身的字段
		for _, msg := range modelMessages {
			if err := tx.Updates(msg).Error; err != nil {
				return gormutils.TranslateError(err)
			}
		}

		return nil
	})
}

func (repo *RoomMessageRepository) FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*domain.RoomMessage, error) {
	var modelMessage model.RoomMessage
	err := repo.db.WithContext(ctx).
		Preload("Recipients").
		Preload("States").
		First(&modelMessage, "id = ?", messageID).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&modelMessage), nil
}

func (repo *RoomMessageRepository) FindsByState(ctx context.Context, userID kernel.UserID, state domain.MessageState) ([]*domain.RoomMessage, error) {
	var modelMessages []*model.RoomMessage
	err := repo.db.WithContext(ctx).
		Joins("JOIN notification_room_message_states ON notification_room_message_states.message_id = notification_room_messages.id").
		Where("notification_room_message_states.user_id = ? AND notification_room_message_states.state = ?", userID, state).
		Preload("Recipients").
		Preload("States").
		Limit(maxLimit).
		Find(&modelMessages).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomains(modelMessages), nil
}

func (repo *RoomMessageRepository) toModel(message *domain.RoomMessage) *model.RoomMessage {
	recipients := make([]*model.RoomMessageRecipient, 0, len(message.RecipientIDs()))
	for _, r := range message.RecipientIDs() {
		recipients = append(recipients, &model.RoomMessageRecipient{
			MessageID: message.ID(),
			UserID:    r,
		})
	}

	states := make([]*model.RoomMessageState, 0, len(message.States()))
	for userID, state := range message.States() {
		states = append(states, &model.RoomMessageState{
			MessageID: message.ID(),
			UserID:    userID,
			State:     state,
		})
	}

	return &model.RoomMessage{
		ID:         message.ID(),
		SenderID:   message.SenderID(),
		RoomID:     message.RoomID(),
		Content:    message.Content(),
		SentAt:     message.SentAt(),
		Recipients: recipients,
		States:     states,
	}
}

func (repo *RoomMessageRepository) toModels(messages []*domain.RoomMessage) []*model.RoomMessage {
	modelMessages := make([]*model.RoomMessage, 0, len(messages))
	for _, roomMessage := range messages {
		modelMessages = append(modelMessages, repo.toModel(roomMessage))
	}
	return modelMessages
}

func (repo *RoomMessageRepository) toDomain(message *model.RoomMessage) *domain.RoomMessage {
	recipients := make([]kernel.UserID, 0, len(message.Recipients))
	for _, r := range message.Recipients {
		recipients = append(recipients, r.UserID)
	}

	states := make(map[kernel.UserID]domain.MessageState, len(message.States))
	for _, s := range message.States {
		states[s.UserID] = s.State
	}
	return domain.LoadRoomMessage(message.ID, message.SenderID, message.RoomID, recipients, states, message.Content, message.SentAt)
}

func (repo *RoomMessageRepository) toDomains(messages []*model.RoomMessage) []*domain.RoomMessage {
	domainMessages := make([]*domain.RoomMessage, 0, len(messages))
	for _, message := range messages {
		domainMessages = append(domainMessages, repo.toDomain(message))
	}
	return domainMessages
}
