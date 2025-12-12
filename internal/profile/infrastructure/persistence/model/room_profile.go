package model

import (
	"gochat/internal/shared/kernel"
	"time"
)

type RoomProfile struct {
	ID           kernel.RoomID `gorm:"primaryKey;type:char(36)"`
	Name         string
	Introduction string
	CreatedAt    time.Time
}

func (*RoomProfile) TableName() string {
	return "profile_room_profiles"
}
