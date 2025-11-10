package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

// TODO 删除对 gorm model的依赖

type Message struct {
	gorm.Model
	ID       domain.MessageID `gorm:"primaryKey;type:char(36)"`
	Content  string           `gorm:"type:text;not null"`
	SenderID kernel.UserID    `gorm:"type:char(36);not null;index"`

	// 用户 <-> 消息 一对多
	Sender *User `gorm:"foreignKey:SenderID;references:ID"`

	// “用户+消息”组合，便于取状态
	States []*MessageState `gorm:"foreignKey:MessageID;references:ID"`
}

func (*Message) TableName() string {
	return "gochat.notification_messages"
}
