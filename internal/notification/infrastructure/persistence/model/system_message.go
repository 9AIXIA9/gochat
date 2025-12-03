package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type SystemMessage struct {
	ID          kernel.MessageID          `gorm:"primaryKey;type:char(36)"`
	Topic       domain.SystemMessageTopic `gorm:"type:varchar(100);index,not null"`
	RecipientID kernel.UserID             `gorm:"type:char(36);not null;index"`
	State       domain.MessageState       `gorm:"type:varchar(36);index,not null"`
	Content     []byte                    `gorm:"type:json;not null"`
	SentAt      time.Time
}

func (*SystemMessage) TableName() string {
	return "notification_system_messages"
}
