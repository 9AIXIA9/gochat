package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID      domain.MessageID `gorm:"primaryKey;type:char(36)"`
	Content string           `gorm:"type:text;not null"`
	SentAt  time.Time        `gorm:"not null"`

	Sender      kernel.UserID `gorm:"not null"`
	RecipientID kernel.UserID `gorm:"foreignKey:RecipientID"`
}

func (m *Message) TableName() string {
	return "gochat.chat_messages"
}
