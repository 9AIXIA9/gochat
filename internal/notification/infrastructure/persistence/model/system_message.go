package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type SystemMessage struct {
	ID          kernel.MessageID    `gorm:"primaryKey;type:char(36)"`
	RecipientID kernel.UserID       `gorm:"type:char(36);not null;index"`
	State       domain.MessageState `gorm:"type:varchar(36);index,not null"`
	Content     string              `gorm:"type:varchar(255);not null"`
	SentAt      time.Time
}

func (*SystemMessage) TableName() string {
	return "notification_system_messages"
}
