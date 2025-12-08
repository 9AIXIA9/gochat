package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
)

type Friendship struct {
	ID      domain.FriendshipID `gorm:"primaryKey;type:char(36)"`
	UserID1 kernel.UserID       `gorm:"type:char(36);not null;uniqueIndex:idx_friendship_pair,priority:1"`
	UserID2 kernel.UserID       `gorm:"type:char(36);not null;uniqueIndex:idx_friendship_pair,priority:2"`

	// 外键到用户，删除受限，更新级联
	_ struct{} `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;foreignKey:UserID1;references:ID"`
	_ struct{} `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;foreignKey:UserID2;references:ID"`
}

func (*Friendship) TableName() string { return "chat_friendships" }
