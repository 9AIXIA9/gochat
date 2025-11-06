package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type MessageUser struct {
	// 复合主键：一个用户-消息的组合
	MessageID domain.MessageID    `gorm:"primaryKey;type:char(36)"`
	UserID    kernel.UserID       `gorm:"primaryKey;type:char(36)"`
	State     domain.MessageState `gorm:"type:varchar(20);not null;index"`
}

func (MessageUser) TableName() string {
	return "gochat.chat_message_users"
}
