package model

import (
	"gorm.io/gorm"
	"log"
)

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&User{},
		&Room{},
		&Chat{},
		&UserRoom{},
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

func (m *Chat) TableName() string {
	return "chat"
}

func (ur *UserRoom) TableName() string {
	return "user_room"
}
