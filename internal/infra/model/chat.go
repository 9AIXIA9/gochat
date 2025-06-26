package model

import (
	"gochat/internal/domain"
	"gorm.io/gorm"
	"time"
)

type Chat struct {
	gorm.Model
	UserNumber domain.UserNumber
	RoomNumber domain.RoomNumber
	Content    string
	SendTime   time.Time
}
