package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type Notice struct {
	gorm.Model
	ID        domain.NoticeID `gorm:"primaryKey;type:char(36)"`
	Theme     domain.NoticeTheme
	Recipient kernel.UserID
	Title     string
	Content   string
	Contact   string
}

func (n *Notice) TableName() string {
	return "gochat.notification_notice"
}
