package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type MessageID kernel.ID

func (i MessageID) String() string {
	return string(i)
}

type MessageType string

const (
	PrivateType MessageType = "private"
)

func (t MessageType) String() string {
	return string(t)
}

type MessageState string

const (
	MessageStateCreated   MessageState = "created"
	MessageStateDelivered MessageState = "delivered"
	//MessageStateRead      MessageState = "read"
)

func (s MessageState) String() string {
	return string(s)
}

type MessageInformation struct {
	id      MessageID
	mType   MessageType
	sender  kernel.UserID
	content string
	sentAt  time.Time
}

func NewMessageInformation(
	id MessageID,
	mType MessageType,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
) *MessageInformation {
	return &MessageInformation{
		id:      id,
		mType:   mType,
		sender:  sender,
		content: content,
		sentAt:  sentAt,
	}
}

func (m *MessageInformation) ID() MessageID {
	return m.id
}

func (m *MessageInformation) Type() MessageType {
	return m.mType
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
