package model

import (
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	ID       kernel.MessageID `gorm:"primaryKey;type:char(36)"`
	SenderID kernel.UserID    `gorm:"type:char(36);not null;index"`
	RoomID   kernel.RoomID    `gorm:"type:char(36);not null;index"`
	Content  string           `gorm:"type:text;not null"`
	SentAt   time.Time

	// 关联收件人
	Recipients []*RoomMessageRecipient `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE"`
	// 关联状态
	States []*RoomMessageState `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE"`
}

func (RoomMessage) TableName() string {
	return "notification_room_messages"
}

type RoomMessageRecipient struct {
	MessageID kernel.MessageID `gorm:"type:char(36);index;not null"`
	UserID    kernel.UserID    `gorm:"type:char(36);index;not null"`
}

func (*RoomMessageRecipient) TableName() string {
	return "notification_room_message_recipients"
}

type RoomMessageState struct {
	MessageID kernel.MessageID    `gorm:"type:char(36);index;not null"`
	UserID    kernel.UserID       `gorm:"type:char(36);index;not null"`
	State     domain.MessageState `gorm:"type:varchar(36);index,not null"`
}

func (*RoomMessageState) TableName() string {
	return "notification_room_message_states"
}
