package model

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type Roomship struct {
	ID        domain.RoomshipID   `gorm:"primaryKey;type:char(36)"`
	RoomID    kernel.RoomID       `gorm:"not null;index:idx_roomship_pair,priority:1"`
	UserID    kernel.UserID       `gorm:"not null;index:idx_roomship_pair,priority:2"`
	Role      domain.RoomshipRole `gorm:"not null;index;type:char(20)"`
	CreatedAt time.Time           `gorm:"not null;index"`
}

func (*Roomship) TableName() string { return "roomship_roomships" }
