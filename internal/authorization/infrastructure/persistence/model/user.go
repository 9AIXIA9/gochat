package model

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type User struct {
	ID                kernel.UserID            `gorm:"primaryKey;type:char(36)"`
	Email             kernel.Email             `gorm:"type:varchar(254);uniqueIndex;not null"`
	Number            kernel.UserNumber        `gorm:"type:varchar(20);uniqueIndex;not null"`
	PasswordEncrypted domain.PasswordEncrypted `gorm:"type:varchar(255);not null"`
	SignedUpAt        time.Time
}

func (u *User) TableName() string {
	return "authorization_users"
}
