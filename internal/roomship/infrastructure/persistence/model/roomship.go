package model

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type Roomship struct {
	ID domain.RoomshipID `gorm:"primaryKey;type:char(36)"`

	RoomID    kernel.RoomID       `gorm:"type:char(36);not null;uniqueIndex:idx_roomship_pair,priority:1"`
	UserID    kernel.UserID       `gorm:"type:char(36);not null;uniqueIndex:idx_roomship_pair,priority:2"`
	Role      domain.RoomshipRole `gorm:"not null;index;type:char(20)"`
	CreatedAt time.Time           `gorm:"not null;index"`

	_ struct{} `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;foreignKey:RoomID;references:ID"`
	_ struct{} `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;foreignKey:UserID;references:ID"`
}

func (*Roomship) TableName() string { return "roomship_roomships" }
