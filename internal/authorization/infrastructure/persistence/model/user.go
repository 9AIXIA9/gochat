package model

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                kernel.UserID     `gorm:"primaryKey;type:char(36)"`
	Email             kernel.Email      `gorm:"type:varchar(254);uniqueIndex;not null"`
	Number            domain.UserNumber `gorm:"type:varchar(20);uniqueIndex;not null"`
	PasswordEncrypted string            `gorm:"type:varchar(255);not null"`
	LastLoggedInAt    time.Time         `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
}

func (u *User) TableName() string {
	return "gochat.authorization_users"
}
