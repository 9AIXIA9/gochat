package model

import (
	"gochat/internal/shared/kernel"
)

type User struct {
	ID     kernel.UserID     `gorm:"primaryKey;type:char(36)"`
	Number kernel.UserNumber `gorm:"type:varchar(20);uniqueIndex;not null"`

	// 双向好友：通过两侧关联汇总
	FriendshipsAsUser1 []*Friendship `gorm:"foreignKey:User1ID;references:ID"`
	FriendshipsAsUser2 []*Friendship `gorm:"foreignKey:User2ID;references:ID"`

	FriendRequestsSent     []*FriendRequest `gorm:"foreignKey:From;references:ID"`
	FriendRequestsReceived []*FriendRequest `gorm:"foreignKey:To;references:ID"`
}

func (*User) TableName() string {
	return "friendship_users"
}
