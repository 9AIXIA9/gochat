package model

import (
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	ID       kernel.MessageID `gorm:"primaryKey;type:char(36)"`
	Content  string           `gorm:"type:text;not null"`
	RoomID   kernel.RoomID    `gorm:"type:char(36);not null;index"`
	SenderID kernel.UserID    `gorm:"type:char(36);not null;index"`

	SentAt time.Time

	Sender *User `gorm:"foreignKey:SenderID;references:ID"`
	Room   *Room `gorm:"foreignKey:RoomID;references:ID"`
}

func (*RoomMessage) TableName() string {
	return "chat_room_messages"
}
