package model

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Number  domain.UserNumber `gorm:"uniqueIndex"`
	PwdHash string
}

type Room struct {
	gorm.Model
	Number       domain.RoomNumber `gorm:"uniqueIndex"`
	Owner        domain.UserNumber
	SecretHash   string
	CurrentUsers int
	MaxUsers     int
}

type Message struct {
	gorm.Model
	ID        domain.MessageID `gorm:"primaryKey"`
	Sender    domain.UserNumber
	Recipient domain.BaseNumber
	Content   string
	SentAt    time.Time
	Type      domain.MessageType
}

type UserRoom struct {
	UserNumber domain.UserNumber `gorm:"primaryKey"`
	RoomNumber domain.RoomNumber `gorm:"primaryKey"`
}

type UserMessage struct {
	gorm.Model
	Sent       bool
	UserNumber domain.UserNumber
	MessageID  domain.MessageID
}
