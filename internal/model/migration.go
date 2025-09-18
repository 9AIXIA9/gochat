package model

import (
	"gorm.io/gorm"
	"log"
)

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&User{},
		&Room{},
		&Message{},
		&UserRoom{},
		&UserMessage{},
	); err != nil {
		log.Fatalf("auto migrate tables failed,err:%v", err)
	}
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
