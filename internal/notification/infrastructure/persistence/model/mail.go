package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type Mail struct {
	gorm.Model
	ID        domain.MailID `gorm:"primaryKey;type:char(36)"`
	Theme     domain.MailTheme
	Recipient kernel.UserID
	Title     string
	Content   string
	Email     kernel.Email
}

func (n *Mail) TableName() string {
	return "gochat.notification_mails"
}
