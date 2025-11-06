package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID     kernel.UserID     `gorm:"primaryKey;type:char(36)"`
	Number domain.UserNumber `gorm:"type:varchar(20);uniqueIndex;not null"`

	// 用户 <-> 房间 多对多
	Rooms []*Room `gorm:"many2many:gochat.chat_room_members;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:RoomID"`

	// 用户 <-> 发送的消息 一对多
	MessageSent []*Message `gorm:"foreignKey:SenderID;references:ID"`

	// “用户+接收消息”组合（用于关联状态）
	MessageLinks []*MessageUser `gorm:"foreignKey:UserID;references:ID"`
}

func (u *User) TableName() string {
	return "gochat.chat_users"
}
