package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type RecipientMessageState struct {
	gorm.Model
	MessageID domain.MessageID `gorm:"uniqueIndex:idx_message_recipient;type:char(36)"`
	Recipient kernel.UserID    `gorm:"uniqueIndex:idx_message_recipient;type:char(36)"`
	State     domain.MessageState

	// Belongs To 关系：状态记录属于一条消息
	Message Message `gorm:"foreignKey:MessageID;references:ID"`
}

func (m *RecipientMessageState) TableName() string {
	return "gochat.chat_recipients_messages_states"
}
