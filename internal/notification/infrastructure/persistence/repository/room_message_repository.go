package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm/clause"
)

var _ domain.RoomMessageRepository = (*RoomMessageRepository)(nil)

type RoomMessageRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewRoomMessageRepository(unitOfWork *gormutils.UnitOfWork) *RoomMessageRepository {
	return &RoomMessageRepository{unitOfWork: unitOfWork}
}

func (repo *RoomMessageRepository) SaveRoomMessage(ctx context.Context, message *domain.RoomMessage) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"sender_id", "room_id", "content", "sent_at"}),
		},
	).Create(repo.toModel(message)).Error)
}

func (repo *RoomMessageRepository) SaveRoomMessages(ctx context.Context, message []*domain.RoomMessage) error {
	models := make([]*model.RoomMessage, 0, len(message))
	for _, msg := range message {
		models = append(models, repo.toModel(msg))
	}
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&models).Error)
}

func (repo *RoomMessageRepository) FindUndeliveredRoomMessages(ctx context.Context, userID kernel.UserID) ([]*domain.RoomMessage, error) {
	var msgs []model.RoomMessage
	if err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Joins("JOIN notification_room_message_recipients AS recipients ON recipients.message_id = notification_room_messages.id").
		Joins("JOIN notification_room_message_states AS states ON states.message_id = notification_room_messages.id AND states.user_id = recipients.user_id").
		Where("recipients.user_id = ? AND states.state = ?", userID, domain.MessageStateUndelivered).
		Preload("Recipients").
		Preload("States").
		Find(&msgs).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}

	messages := make([]*domain.RoomMessage, 0, len(msgs))
	for _, m := range msgs {
		messages = append(messages, repo.toDomain(&m))
	}
	return messages, nil
}

func (repo *RoomMessageRepository) FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*domain.RoomMessage, error) {
	var message model.RoomMessage

	if err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Preload("Recipients").
		Preload("States").
		Where("id = ?", messageID).First(&message).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}

	return repo.toDomain(&message), nil
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
