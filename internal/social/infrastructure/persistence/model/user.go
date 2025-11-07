package model

import "gochat/internal/shared/kernel"

type User struct {
	ID kernel.UserID `gorm:"primaryKey;type:char(36)"`

	Rooms []*Room `gorm:"many2many:gochat.social_room_members;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:RoomID"`
}

func (u *User) TableName() string {
	return "gochat.social_users"
}
