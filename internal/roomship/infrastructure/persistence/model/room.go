package model

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type Room struct {
	ID                kernel.RoomID            `gorm:"primaryKey;type:char(36)"`
	Number            domain.RoomNumber        `gorm:"type:varchar(20);uniqueIndex;not null"`
	OwnerID           kernel.UserID            `gorm:"not null;index"`
	PasswordEncrypted domain.PasswordEncrypted `gorm:"type:varchar(255)"`
	MaxMemberCount    int                      `gorm:"not null;default:0;check:max_member_count >= 0"`
	CreatedAt         time.Time
}

func (*Room) TableName() string { return "roomship_rooms" }
