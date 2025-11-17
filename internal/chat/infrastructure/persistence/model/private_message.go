package model

import (
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessage struct {
	ID          kernel.MessageID `gorm:"primaryKey;type:char(36)"`
	Content     string           `gorm:"type:text;not null"`
	RecipientID kernel.UserID    `gorm:"type:char(36);not null;index"`
	SenderID    kernel.UserID    `gorm:"type:char(36);not null;index"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Sender    *User `gorm:"foreignKey:SenderID;references:ID"`
	Recipient *User `gorm:"foreignKey:RecipientID;references:ID"`
}

func (*PrivateMessage) TableName() string {
	return "gochat.chat_private_messages"
}
