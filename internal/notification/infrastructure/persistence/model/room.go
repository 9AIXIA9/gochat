package model

import (
	"gochat/internal/notification/domain"
)

type Room struct {
	ID domain.RoomID `gorm:"primaryKey;type:char(36)"`
}

func (*Room) TableName() string {
	return "gochat.notification_rooms"
}
