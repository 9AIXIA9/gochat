package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type Mail struct {
	gorm.Model
	ID        domain.MailID    `gorm:"primaryKey;type:char(36)"`
	Theme     domain.MailTheme `gorm:"type:varchar(50);not null;index"`
	Recipient kernel.UserID    `gorm:"type:char(36);not null;index"`
	Title     string           `gorm:"type:varchar(200);not null"`
	Content   string           `gorm:"type:text;not null"`
	Email     kernel.Email     `gorm:"type:varchar(254);not null;index"`
}

func (n *Mail) TableName() string {
	return "gochat.notification_mails"
}
