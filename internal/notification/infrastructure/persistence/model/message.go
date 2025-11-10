package model

import "gochat/internal/notification/domain"

type Message struct {
	ID domain.MessageID `gorm:"primaryKey;type:char(36)"`

	// “用户+消息”组合，便于取状态
	States []*MessageState `gorm:"foreignKey:MessageID;references:ID"`
}

func (*Message) TableName() string {
	return "gochat.notification_messages"
}
