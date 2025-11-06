package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID     kernel.UserID     `gorm:"primaryKey;type:char(36)"`
	Number domain.UserNumber `gorm:"type:varchar(20);uniqueIndex"`

	// 发送的消息（一对多）
	MessagesReceived []*Message `gorm:"foreignKey:RecipientID"`

	// 加入的房间（多对多）
	Rooms []Room `gorm:"many2many:room_members;"`
}

func (u *User) TableName() string {
	return "gochat.chat_users"
}
