package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"
)

var _ domain.RoomMessageRepository = (*RoomMessageRepository)(nil)

type RoomMessageRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewRoomMessageRepository(unitOfWork *gormutils.UnitOfWork) *RoomMessageRepository {
	return &RoomMessageRepository{unitOfWork: unitOfWork}
}

func (repo *RoomMessageRepository) SaveRoomMessage(ctx context.Context, message *domain.RoomMessage) error {
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

	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Save(&model.RoomMessage{
		ID:         message.ID(),
		SenderID:   message.SenderID(),
		RoomID:     message.RoomID(),
		Content:    message.Content(),
		SentAt:     message.SentAt(),
		Recipients: recipients,
		States:     states,
	}).Error)
}

func (repo *RoomMessageRepository) FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*domain.RoomMessage, error) {
	var message model.RoomMessage

	if err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Preload("Recipients").
		Preload("States").
		Where("id = ?", messageID).First(&message).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}

	recipients := make([]kernel.UserID, 0, len(message.Recipients))
	for _, r := range message.Recipients {
		recipients = append(recipients, r.UserID)
	}

	states := make(map[kernel.UserID]domain.MessageState, len(message.States))
	for _, s := range message.States {
		states[s.UserID] = s.State
	}
	return domain.LoadRoomMessage(message.ID, message.SenderID, message.RoomID, recipients, states, message.Content, message.SentAt), nil
}

func (repo *RoomMessageRepository) FindReceivedRoomMessage(ctx context.Context, userID kernel.UserID) ([]*domain.RoomMessage, error) {
	var messages []model.RoomMessage

	if err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Joins("JOIN notification_room_message_recipients r ON r.message_id = notification_room_messages.id").
		Where("r.user_id = ?", userID).
		Preload("Recipients").
		Preload("States").
		Find(&messages).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}

	domainMessages := make([]*domain.RoomMessage, 0, len(messages))
	for _, m := range messages {
		recipients := make([]kernel.UserID, 0, len(m.Recipients))
		for _, r := range m.Recipients {
			recipients = append(recipients, r.UserID)
		}

		states := make(map[kernel.UserID]domain.MessageState, len(m.States))
		for _, s := range m.States {
			states[s.UserID] = s.State
		}
		domainMessages = append(domainMessages, domain.LoadRoomMessage(
			m.ID,
			m.SenderID,
			m.RoomID,
			recipients,
			states,
			m.Content,
			m.SentAt,
		))
	}
	return domainMessages, nil
}
