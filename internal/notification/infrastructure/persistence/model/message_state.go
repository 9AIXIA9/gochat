package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type MessageState struct {
	// 复合主键：一个用户-消息的组合
	MessageID domain.MessageID    `gorm:"primaryKey;type:char(36)"`
	UserID    kernel.UserID       `gorm:"primaryKey;type:char(36)"`
	State     domain.MessageState `gorm:"type:varchar(20);not null;index"`

	// 关联
	Message *Message `gorm:"foreignKey:MessageID;references:ID"`
}

func (*MessageState) TableName() string {
	return "gochat.notification_message_states"
}
