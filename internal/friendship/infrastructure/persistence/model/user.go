package model

import (
	"gochat/internal/shared/kernel"
)

type User struct {
	ID     kernel.UserID     `gorm:"primaryKey;type:char(36)"`
	Number kernel.UserNumber `gorm:"type:varchar(20);uniqueIndex;not null"`
}

func (*User) TableName() string {
	return "friendship_users"
}
