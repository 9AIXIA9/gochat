package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID        domain.MessageID
	Type      domain.MessageType
	Recipient kernel.ID
	State     domain.MessageState
	Sender    kernel.UserID
	Content   string
	SentAt    time.Time
}

func (m *Message) TableName() string {
	return "gochat.chat_messages"
}
