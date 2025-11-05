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
	State     domain.State
}

func (m *RecipientMessageState) TableName() string {
	return "gochat.chat_recipients_messages_states"
}
