package model

import (
	"gochat/internal/shared/kernel"
)

type Room struct {
	ID     kernel.RoomID     `gorm:"primaryKey;type:char(36)"`
	Number kernel.RoomNumber `gorm:"type:varchar(20);uniqueIndex;not null"`

	// 房间 <-> 用户 多对多
	Members []*User `gorm:"many2many:gochat.chat_room_members;foreignKey:ID;joinForeignKey:RoomID;references:ID;joinReferences:UserID"`

	// 房间 <-> 收到的消息 一对多
	MessagesReceived []*RoomMessage `gorm:"foreignKey:RoomID;references:ID"`
}

func (*Room) TableName() string {
	return "chat_rooms"
}
