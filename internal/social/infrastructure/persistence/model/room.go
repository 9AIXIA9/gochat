package model

import (
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
	"time"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	ID                domain.RoomID `gorm:"primaryKey;type:char(36)"`
	Owner             kernel.UserID
	Number            domain.RoomNumber
	PasswordEncrypted string
	MemberCount       int
	MaxMemberCount    int
	CreatedAt         time.Time
	Members           []*RoomMember `gorm:"foreignKey:RoomID;references:ID"`
}

func (r *Room) TableName() string {
	return "gochat.social_rooms"
}
