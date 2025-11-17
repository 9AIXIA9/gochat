package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type Message struct {
	id      kernel.MessageID
	sender  kernel.UserID
	content string
	sentAt  time.Time
}

func NewMessage(id kernel.MessageID, sender kernel.UserID, content string, sentAt time.Time) *Message {
	return &Message{id: id, sender: sender, content: content, sentAt: sentAt}
}

func (m *Message) ID() kernel.MessageID {
	return m.id
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
