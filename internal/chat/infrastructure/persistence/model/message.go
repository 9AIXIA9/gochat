package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID      domain.MessageID `gorm:"primaryKey;type:char(36)"`
	Content string           `gorm:"type:text;not null"`

	SenderID kernel.UserID `gorm:"type:char(36);not null;index"`
	Sender   *User         `gorm:"foreignKey:SenderID;references:ID"`

	// “用户+消息”组合，便于取状态
	UserStates []*UserMessageState `gorm:"foreignKey:MessageID;references:ID"`
}

func (m *Message) TableName() string {
	return "gochat.chat_messages"
}
