package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type Message struct {
	ID       domain.MessageID `gorm:"primaryKey;type:char(36)"`
	Content  string           `gorm:"type:text;not null"`
	SenderID kernel.UserID    `gorm:"type:char(36);not null;index"`

	CreatedAt time.Time

	// “用户+消息”组合，便于取状态
	States []*MessageState `gorm:"foreignKey:MessageID;references:ID"`
}

func (*Message) TableName() string {
	return "gochat.notification_messages"
}
