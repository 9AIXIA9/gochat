package model

import (
	"gochat/internal/shared/kernel"
)

type User struct {
	ID kernel.UserID `gorm:"primaryKey;type:char(36)"`

	// 用户 <-> 房间 多对多
	Rooms []*Room `gorm:"many2many:gochat.notification_room_members;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:RoomID"`

	// 用户 <-> 消息 一对多
	MessageSent []*Message `gorm:"foreignKey:SenderID;references:ID"`

	// “用户+接收消息”组合（用于关联状态）
	MessageStates []*MessageState `gorm:"foreignKey:UserID;references:ID"`
}

func (*User) TableName() string {
	return "gochat.notification_users"
}
