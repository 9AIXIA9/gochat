package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type Message struct {
	sender  kernel.UserID
	content string
	sentAt  time.Time
}

func NewMessage(sender kernel.UserID, content string, sentAt time.Time) *Message {
	return &Message{sender: sender, content: content, sentAt: sentAt}
}

func (m *Message) Sender() kernel.UserID {
	return m.sender
}

func (m *Message) Content() string {
	return m.content
}

func (m *Message) SentAt() time.Time {
	return m.sentAt
}
