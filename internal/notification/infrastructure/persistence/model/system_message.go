package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type SystemMessage struct {
	ID          kernel.MessageID    `gorm:"primaryKey;type:char(36)"`
	RecipientID kernel.UserID       `gorm:"type:char(36);not null;index"`
	State       domain.MessageState `gorm:"type:varchar(36);not null;index"`
	Content     string              `gorm:"type:varchar(255);not null"`
	SentAt      time.Time

	_ struct{} `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;foreignKey:RecipientID;references:ID"`
}

func (*SystemMessage) TableName() string {
	return "notification_system_messages"
}
