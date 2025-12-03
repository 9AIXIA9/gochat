package model

import (
	"gochat/internal/shared/kernel"
)

type Room struct {
	ID kernel.RoomID `gorm:"primaryKey;type:char(36)"`

	// 房间 <-> 用户 多对多
	Members []*User `gorm:"many2many:gochat.chat_room_members;foreignKey:ID;joinForeignKey:RoomID;references:ID;joinReferences:UserID"`
}

func (*Room) TableName() string {
	return "chat_rooms"
}
