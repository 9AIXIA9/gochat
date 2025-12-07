package model

import (
	"gochat/internal/shared/kernel"
)

type Room struct {
	ID kernel.RoomID `gorm:"primaryKey;type:char(36)"`
}

func (*Room) TableName() string {
	return "chat_rooms"
}
