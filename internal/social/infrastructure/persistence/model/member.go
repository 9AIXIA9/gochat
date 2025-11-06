package model

import (
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"

	"gorm.io/gorm"
)

type RoomMember struct {
	gorm.Model
	RoomID domain.RoomID `gorm:"uniqueIndex:idx_room_member;type:char(36);not null"`
	Member kernel.UserID `gorm:"uniqueIndex:idx_room_member;type:char(36);not null"`

	Room *Room `gorm:"foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *RoomMember) TableName() string {
	return "gochat.social_room_members"
}
