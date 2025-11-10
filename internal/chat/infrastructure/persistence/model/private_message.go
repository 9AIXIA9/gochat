package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type PrivateMessage struct {
	gorm.Model
	ID          domain.MessageID `gorm:"primaryKey;type:char(36)"`
	Content     string           `gorm:"type:text;not null"`
	RecipientID kernel.UserID    `gorm:"type:char(36);not null;index"`
	SenderID    kernel.UserID    `gorm:"type:char(36);not null;index"`

	Sender    *User `gorm:"foreignKey:SenderID;references:ID"`
	Recipient *User `gorm:"foreignKey:RecipientID;references:ID"`
}

func (*PrivateMessage) TableName() string {
	return "gochat.chat_private_messages"
}
