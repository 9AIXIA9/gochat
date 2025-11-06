package model

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID      domain.MessageID `gorm:"primaryKey;type:char(36)"`
	Content string           `gorm:"type:text;not null"`

	Sender      kernel.UserID `gorm:"not null;index"`
	RecipientID kernel.UserID `gorm:"not null;index"`

	// 关联到发送者与接收者
	SenderUser *User `gorm:"foreignKey:Sender;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Recipient  *User `gorm:"foreignKey:RecipientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (m *Message) TableName() string {
	return "gochat.chat_messages"
}
