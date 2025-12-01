package model

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type FriendRequest struct {
	ID      kernel.OperationID        `gorm:"primaryKey;type:char(36)"`
	From    kernel.UserID             `gorm:"type:char(36);not null;index:idx_from_to,unique"`
	To      kernel.UserID             `gorm:"type:char(36);not null;index:idx_from_to,unique"`
	Content string                    `gorm:"type:varchar(255);not null"`
	State   domain.FriendRequestState `gorm:"type:varchar(20);not null"`
	SentAt  time.Time                 //UTC

	FromUser *User `gorm:"foreignKey:From;references:ID"`
	ToUser   *User `gorm:"foreignKey:To;references:ID"`
}

func (*FriendRequest) TableName() string {
	return "friendship_friend_requests"
}
