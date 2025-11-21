package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
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
	return gormutils.TranslateError(repo.unitOfWork.DB(ctx).WithContext(ctx).Create(&model.RoomMessage{
		ID:       message.ID(),
		Content:  message.Content(),
		RoomID:   message.RoomID(),
		SenderID: message.SenderID(),
		SentAt:   message.SentAt(),
	}).Error)
}

func (repo *RoomMessageRepository) FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*domain.RoomMessage, error) {
	var m model.RoomMessage
	err := repo.unitOfWork.DB(ctx).WithContext(ctx).
		Where("id = ?", messageID).
		First(&m).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return domain.LoadRoomMessage(m.ID, m.SenderID, m.RoomID, m.Content, m.SentAt), nil
}
