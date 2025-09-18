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

type Message struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	Sender    domain.UserNumber
	Recipient domain.BaseNumber
	Content   string
	SentAt    time.Time
	Sent      bool
}

type UserRoom struct {
	UserNumber domain.UserNumber `gorm:"primaryKey"`
	RoomNumber domain.RoomNumber `gorm:"primaryKey"`
}
