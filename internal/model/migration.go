package model

import (
	"gorm.io/gorm"
	"log"
)

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&GormUser{},
		&GormRoom{},
		&GormMessage{},
		&GormUserRoom{},
	); err != nil {
		log.Fatalf("auto migrate tables failed,err:%v", err)
	}
}

func (u *GormUser) TableName() string {
	return "user"
}

func (r *GormRoom) TableName() string {
	return "room"
}

func (m *GormMessage) TableName() string {
	return "message"
}

func (ur *GormUserRoom) TableName() string {
	return "user_room"
}
