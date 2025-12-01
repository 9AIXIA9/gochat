package model

import (
	"fmt"
	"gochat/internal/shared/kernel"
	"time"

	"gorm.io/gorm"
)

type Friendship struct {
	User1ID   kernel.UserID `gorm:"column:user1_id;type:char(36);not null;index:idx_user_pair,unique"`
	User2ID   kernel.UserID `gorm:"column:user2_id;type:char(36);not null;index:idx_user_pair,unique"`
	CreatedAt time.Time

	User1 *User `gorm:"foreignKey:User1ID;references:ID;constraint:OnDelete:CASCADE"`
	User2 *User `gorm:"foreignKey:User2ID;references:ID;constraint:OnDelete:CASCADE"`
}

func (*Friendship) TableName() string { return "friendship_friendships" }

func (f *Friendship) BeforeCreate(*gorm.DB) error {
	if f.User1ID == f.User2ID {
		return fmt.Errorf("cannot create friendship with self")
	}
	if string(f.User1ID) > string(f.User2ID) {
		f.User1ID, f.User2ID = f.User2ID, f.User1ID
	}
	return nil
}
