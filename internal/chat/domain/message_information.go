package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type MessageInformation struct {
	sender  kernel.UserID
	content string
	sentAt  time.Time
}

func NewMessageInformation(
	sender kernel.UserID,
	content string,
	sentAt time.Time,
) *MessageInformation {
	return &MessageInformation{
		sender:  sender,
		content: content,
		sentAt:  sentAt,
	}
}

func (m *MessageInformation) Sender() kernel.UserID {
	return m.sender
}

func (m *MessageInformation) Content() string {
	return m.content
}

func (m *MessageInformation) SentAt() time.Time {
	return m.sentAt
}
