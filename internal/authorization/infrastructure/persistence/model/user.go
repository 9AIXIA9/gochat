package model

import (
	"gochat/internal/shared/kernel"
	"time"
)

type User struct {
	ID                kernel.UserID     `gorm:"primaryKey;type:char(36)"`
	Email             kernel.Email      `gorm:"type:varchar(254);uniqueIndex;not null"`
	Number            kernel.UserNumber `gorm:"type:varchar(20);uniqueIndex;not null"`
	PasswordEncrypted string            `gorm:"type:varchar(255);not null"`

	LastLoggedInAt time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
	CreatedAt      time.Time
}

func (u *User) TableName() string {
	return "authorization_users"
}
