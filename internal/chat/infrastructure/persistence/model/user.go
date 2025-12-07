package model

import (
	"gochat/internal/shared/kernel"
)

type User struct {
	ID kernel.UserID `gorm:"primaryKey;type:char(36)"`
}

func (*User) TableName() string {
	return "chat_users"
}
