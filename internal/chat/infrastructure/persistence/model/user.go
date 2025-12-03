package model

import (
	"gochat/internal/shared/kernel"
)

type User struct {
	ID kernel.UserID `gorm:"primaryKey;type:char(36)"`

	// 用户 <-> 房间 多对多
	Rooms []*Room `gorm:"many2many:gochat.chat_room_members;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:RoomID"`

	// 用户 <-> 发送的消息 一对多
	PrivateMessagesSent     []*PrivateMessage `gorm:"foreignKey:SenderID;references:ID"`
	RoomMessagesSent        []*RoomMessage    `gorm:"foreignKey:SenderID;references:ID"`
	PrivateMessagesReceived []*PrivateMessage `gorm:"foreignKey:RecipientID;references:ID"`
}

func (*User) TableName() string {
	return "chat_users"
}
