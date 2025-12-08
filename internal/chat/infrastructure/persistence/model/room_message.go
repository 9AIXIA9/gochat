package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	ID       kernel.MessageID `gorm:"primaryKey;type:char(36)"`
	SenderID kernel.UserID    `gorm:"type:char(36);not null;index"`
	RoomID   kernel.RoomID    `gorm:"type:char(36);not null;index"`
	Content  string           `gorm:"type:text;not null"`
	SentAt   time.Time

	// 外键约束
	_ struct{} `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;foreignKey:SenderID;references:ID"`
	_ struct{} `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;foreignKey:RoomID;references:ID"`

	// 关联状态
	States []*RoomMessageState `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE"`
}

func (RoomMessage) TableName() string {
	return "chat_room_messages"
}

type RoomMessageState struct {
	MessageID kernel.MessageID    `gorm:"type:char(36);not null;index;uniqueIndex:idx_room_msg_state,priority:1"`
	UserID    kernel.UserID       `gorm:"type:char(36);not null;index;uniqueIndex:idx_room_msg_state,priority:2"`
	State     domain.MessageState `gorm:"type:varchar(36);not null;index"`

	_ struct{} `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;foreignKey:MessageID;references:ID"`
	_ struct{} `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;foreignKey:UserID;references:ID"`
}

func (*RoomMessageState) TableName() string {
	return "chat_room_message_states"
}
