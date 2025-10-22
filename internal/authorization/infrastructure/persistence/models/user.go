package models

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID             kernel.UserID     `gorm:"primaryKey;type:char(36)"`
	Email          kernel.Email      `gorm:"type:varchar(254);uniqueIndex"`
	Number         domain.UserNumber `gorm:"type:varchar(20);uniqueIndex"`
	PasswordHash   string
	SignedUpAt     time.Time
	LastLoggedInAt time.Time
}

func (u *User) TableName() string {
	return "gochat.authorization_user"
}
