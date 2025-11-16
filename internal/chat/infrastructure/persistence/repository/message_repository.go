package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
)

var _ application.MessageRepository = (*MessageRepository)(nil)

type MessageRepository struct {
	unitOfWork *gormutils.UnitOfWork
}

func NewMessageRepository(unitOfWork *gormutils.UnitOfWork) *MessageRepository {
	return &MessageRepository{unitOfWork: unitOfWork}
}

func (repo *MessageRepository) SavePrivateMessage(ctx context.Context, recipient kernel.UserID, message *domain.Message) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&model.PrivateMessage{
		ID:          message.ID(),
		Content:     message.Content(),
		RecipientID: recipient,
		SenderID:    message.Sender(),
		CreatedAt:   message.SentAt(),
	}).Error)
}

func (repo *MessageRepository) SaveRoomMessage(ctx context.Context, roomID domain.RoomID, message *domain.Message) error {
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&model.RoomMessage{
		ID:       message.ID(),
		Content:  message.Content(),
		RoomID:   roomID,
		SenderID: message.Sender(),
	}).Error)
}

func (repo *MessageRepository) FindPrivateMessage(ctx context.Context, recipient kernel.UserID, messageID domain.MessageID) (*domain.Message, error) {
	var m model.PrivateMessage
	err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Where("id = ? AND recipient_id = ?", messageID, recipient).
		First(&m).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.NewMessage(m.ID, m.SenderID, m.Content, m.CreatedAt), nil
}

func (repo *MessageRepository) FindRoomMessage(ctx context.Context, roomID domain.RoomID, messageID domain.MessageID) (*domain.Message, error) {
	var m model.RoomMessage
	err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Where("id = ? AND room_id = ?", messageID, roomID).
		First(&m).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.NewMessage(m.ID, m.SenderID, m.Content, m.CreatedAt), nil
}
