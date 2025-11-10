package model

import (
	"gochat/internal/shared/kernel"
)

type User struct {
	ID kernel.UserID `gorm:"primaryKey;type:char(36)"`
}

func (u *User) TableName() string {
	return "gochat.notification_users"
}
