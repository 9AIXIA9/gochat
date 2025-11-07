package model

import (
	"gochat/internal/notification/domain"
)

type Room struct {
	ID      domain.RoomID `gorm:"primaryKey;type:char(36)"`
	Members []*User       `gorm:"many2many:gochat.notification_room_members;foreignKey:ID;joinForeignKey:RoomID;references:ID;joinReferences:UserID"`
}

func (r *Room) TableName() string {
	return "gochat.notification_rooms"
}
