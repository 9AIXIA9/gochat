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
	Sender  kernel.UserID
	Content string
	SentAt  time.Time
	// 一对多关系：一条消息对应多个接收状态
	States []*RecipientMessageState `gorm:"foreignKey:MessageID;references:ID"`
}

func (m *Message) TableName() string {
	return "gochat.chat_messages"
}
