package model

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Number  domain.UserNumber `gorm:"uniqueIndex"`
	Name    string
	PwdHash string
}

type Room struct {
	gorm.Model
	Number       domain.RoomNumber `gorm:"uniqueIndex"`
	Name         string
	Owner        domain.UserNumber
	SecretHash   string
	Description  string
	CurrentUsers int
	MaxUsers     int
}

type Chat struct {
	gorm.Model
	UserNumber domain.UserNumber
	RoomNumber domain.RoomNumber
	Content    string
	SentAt     time.Time
}

type UserRoom struct {
	UserNumber domain.UserNumber `gorm:"primaryKey"`
	RoomNumber domain.RoomNumber `gorm:"primaryKey"`
}
