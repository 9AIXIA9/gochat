package model

import (
	"gochat/internal/shared/kernel"
)

type User struct {
	ID    kernel.UserID `gorm:"primaryKey;type:char(36)"`
	Email kernel.Email  `gorm:"type:varchar(254);uniqueIndex;not null"`
}

func (u *User) TableName() string {
	return "gochat.notification_users"
}
