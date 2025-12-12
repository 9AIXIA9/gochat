package model

import (
	"gochat/internal/shared/kernel"
	"time"
)

type UserProfile struct {
	ID          kernel.UserID `gorm:"primaryKey;type:char(36)"`
	Name        string
	Gender      kernel.Gender
	Email       kernel.Email
	PhoneNumber kernel.PhoneNumber
	Address     kernel.Address
	Sign        string
	SignedUpAt  time.Time
}

func (*UserProfile) TableName() string {
	return "profile_user_profiles"
}
