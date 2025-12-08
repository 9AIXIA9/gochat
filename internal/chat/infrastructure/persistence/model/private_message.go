package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessage struct {
	ID          kernel.MessageID    `gorm:"primaryKey;type:char(36)"`
	Content     string              `gorm:"type:text;not null"`
	RecipientID kernel.UserID       `gorm:"type:char(36);not null;index"`
	SenderID    kernel.UserID       `gorm:"type:char(36);not null;index"`
	State       domain.MessageState `gorm:"type:varchar(36);not null;index"`
	SentAt      time.Time

	_ struct{} `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;foreignKey:SenderID;references:ID"`
	_ struct{} `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;foreignKey:RecipientID;references:ID"`
}

func (*PrivateMessage) TableName() string {
	return "chat_private_messages"
}
