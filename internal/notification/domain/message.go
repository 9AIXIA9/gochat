package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type MessageID kernel.ID

func (i MessageID) String() string {
	return string(i)
}

type Message struct {
	id        MessageID
	sender    kernel.UserID
	recipient kernel.UserID
	content   string
	sentAt    time.Time
}

func NewMessage(
	id MessageID,
	recipient kernel.UserID,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
) *Message {
	return &Message{
		id:        id,
		recipient: recipient,
		sender:    sender,
		content:   content,
		sentAt:    sentAt,
	}
}

func (m *Message) ID() MessageID {
	return m.id
}

func (m *Message) Recipient() kernel.UserID {
	return m.recipient
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
