package model

import (
	"gochat/internal/chat/domain"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	ID     domain.RoomID     `gorm:"primaryKey;type:char(36)"`
	Number domain.RoomNumber `gorm:"type:varchar(20);uniqueIndex"`

	// 房间成员（多对多）
	Members []*User `gorm:"many2many:room_members;"`
}

func (r *Room) TableName() string {
	return "gochat.chat_rooms"
}
