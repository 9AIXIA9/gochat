package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type RoomMessage struct {
	gorm.Model
	ID       domain.MessageID `gorm:"primaryKey;type:char(36)"`
	Content  string           `gorm:"type:text;not null"`
	RoomID   domain.RoomID    `gorm:"type:char(36);not null;index"`
	SenderID kernel.UserID    `gorm:"type:char(36);not null;index"`

	Sender *User `gorm:"foreignKey:SenderID;references:ID"`
	Room   *Room `gorm:"foreignKey:RoomID;references:ID"`
}

func (*RoomMessage) TableName() string {
	return "gochat.chat_room_messages"
}
