package model

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type Friendship struct {
	ID        domain.FriendshipID `gorm:"primaryKey;type:char(36)"`
	UserID1   kernel.UserID       `gorm:"not null;index:idx_friendship_pair,priority:1"`
	UserID2   kernel.UserID       `gorm:"not null;index:idx_friendship_pair,priority:2"`
	CreatedAt time.Time           `gorm:"not null;index"`
}

func (*Friendship) TableName() string { return "friendship_friendships" }
