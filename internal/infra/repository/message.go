package repository

import (
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

func (c *MessageRepository) Create(message domain.Message) (bool, error) {
	gormMessage := model.MessageFromDomain(message)
	if err := c.db.Create(gormMessage).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (c *MessageRepository) QueryAllByRoomNumber(number domain.RoomNumber) ([]domain.Message, error) {
	var gormMessages []model.GormMessage
	if err := c.db.Where("room_number = ?", number).Order("sent_at").Find(&gormMessages).Error; err != nil {
		return nil, err
	}

	messages := make([]domain.Message, len(gormMessages))
	for i, gormMessage := range gormMessages {
		messages[i] = gormMessage.ToDomain()
	}

	return messages, nil
}
