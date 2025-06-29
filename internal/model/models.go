package model

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
	"time"
)

type GormUser struct {
	gorm.Model
	Number  domain.UserNumber `gorm:"uniqueIndex"`
	Name    string
	PwdHash string
}

type GormRoom struct {
	gorm.Model
	Number       domain.RoomNumber `gorm:"uniqueIndex"`
	Name         string
	Owner        domain.UserNumber
	SecretHash   string
	Description  string
	CurrentUsers int
	MaxUsers     int
}

type GormMessage struct {
	gorm.Model
	UserNumber domain.UserNumber
	RoomNumber domain.RoomNumber
	Content    string
	SentAt     time.Time
}

type GormUserRoom struct {
	UserNumber domain.UserNumber `gorm:"primaryKey"`
	RoomNumber domain.RoomNumber `gorm:"primaryKey"`
}
