package model

import (
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
	"time"
)

type Room struct {
	ID                domain.RoomID     `gorm:"primaryKey;type:char(36)"`
	Owner             kernel.UserID     `gorm:"not null;index"`
	Number            domain.RoomNumber `gorm:"type:varchar(20);uniqueIndex;not null"`
	PasswordEncrypted string            `gorm:"type:varchar(255);not null"`
	MemberCount       int               `gorm:"not null;default:0;check:member_count >= 0"`
	MaxMemberCount    int               `gorm:"not null;default:0;check:max_member_count >= 0"`

	CreatedAt time.Time

	Members []*User `gorm:"many2many:gochat.social_room_members;foreignKey:ID;joinForeignKey:RoomID;references:ID;joinReferences:UserID"`
}

func (r *Room) TableName() string {
	return "gochat.social_rooms"
}
