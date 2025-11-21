package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
)

var _ domain.MessageRepository = (*MessageRepository)(nil)

type MessageRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewMessageRepository(unitOfWork *gormutils.UnitOfWork) *MessageRepository {
	return &MessageRepository{unitOfWork: unitOfWork}
}

func (repo *MessageRepository) SavePrivateMessage(ctx context.Context, message *domain.PrivateMessage) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&model.PrivateMessage{
		ID:          message.ID(),
		Content:     message.Content(),
		RecipientID: message.RecipientID(),
		SenderID:    message.SenderID(),
		SentAt:      message.SentAt(),
	}).Error)
}

func (repo *MessageRepository) SaveRoomMessage(ctx context.Context, message *domain.RoomMessage) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&model.RoomMessage{
		ID:       message.ID(),
		Content:  message.Content(),
		RoomID:   message.RoomID(),
		SenderID: message.SenderID(),
		SentAt:   message.SentAt(),
	}).Error)
}

func (repo *MessageRepository) FindPrivateMessage(ctx context.Context, messageID kernel.MessageID) (*domain.PrivateMessage, error) {
	var m model.PrivateMessage
	err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Where("id = ?", messageID).
		First(&m).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.LoadPrivateMessage(m.ID, m.SenderID, m.RecipientID, m.Content, m.SentAt), nil
}

func (repo *MessageRepository) FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*domain.RoomMessage, error) {
	var m model.RoomMessage
	err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Where("id = ?", messageID).
		First(&m).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.LoadRoomMessage(m.ID, m.SenderID, m.RoomID, m.Content, m.SentAt), nil
}
