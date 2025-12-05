package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type Friendship struct {
	ID      domain.FriendshipID `gorm:"primaryKey;type:char(36)"`
	UserID1 kernel.UserID       `gorm:"not null;index:idx_friendship_pair,priority:1"`
	UserID2 kernel.UserID       `gorm:"not null;index:idx_friendship_pair,priority:2"`
}

func (*Friendship) TableName() string { return "chat_friendships" }
