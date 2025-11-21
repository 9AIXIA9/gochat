package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type PrivateMessage struct {
	ID          kernel.MessageID    `gorm:"primaryKey;type:char(36)"`
	SenderID    kernel.UserID       `gorm:"type:char(36);not null;index"`
	RecipientID kernel.UserID       `gorm:"type:char(36);not null;index"`
	State       domain.MessageState `gorm:"type:varchar(10);index,not null"`
	Content     string              `gorm:"type:text;not null"`
	SentAt      time.Time
}

func (*PrivateMessage) TableName() string {
	return "notification_private_messages"
}
