package model

import (
	"gochat/internal/social/domain"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	ID                domain.RoomID     `gorm:"primaryKey;type:char(36)"`
	Number            domain.RoomNumber `gorm:"type:varchar(20);uniqueIndex"`
	PasswordEncrypted string
}

func (r *Room) TableName() string {
	return "gochat.social_rooms"
}
