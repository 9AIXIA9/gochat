package model

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Room{},
		&Message{},
		&UserRoom{},
		&UserMessage{})
}

func (u *User) TableName() string {
	return "user"
}

func (r *Room) TableName() string {
	return "room"
}

func (m *Message) TableName() string {
	return "message"
}

func (um *UserMessage) TableName() string {
	return "user_message"
}
func (ur *UserRoom) TableName() string {
	return "user_room"
}
