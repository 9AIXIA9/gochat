package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"time"

	"gorm.io/gorm"
)

type MessageInformation struct {
	gorm.Model
	ID      domain.MessageID `gorm:"primaryKey;type:char(36)"`
	Sender  kernel.UserID
	Content string
	SentAt  time.Time
}

func (m *MessageInformation) TableName() string {
	return "gochat.chat_messages_information"
}
