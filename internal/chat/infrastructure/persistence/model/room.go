package model

import (
	"gochat/internal/chat/domain"
)

type Room struct {
	ID     domain.RoomID     `gorm:"primaryKey;type:char(36)"`
	Number domain.RoomNumber `gorm:"type:varchar(20);uniqueIndex;not null"`

	// 房间 <-> 用户 多对多
	Members []*User `gorm:"many2many:gochat.chat_room_members;foreignKey:ID;joinForeignKey:RoomID;references:ID;joinReferences:UserID"`

	// 房间 <-> 收到的消息 一对多
	MessagesReceived []*RoomMessage `gorm:"foreignKey:RoomID;references:ID"`
}

func (*Room) TableName() string {
	return "gochat.chat_rooms"
}
