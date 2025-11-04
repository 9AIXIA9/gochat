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
	MessageTypePrivate MessageType = "private"
)

func (t MessageType) String() string {
	return string(t)
}

type MessageState string

const (
	MessageStateSent      MessageState = "sent"
	MessageStateDelivered MessageState = "delivered"
	MessageStateRead      MessageState = "read"
)

func (s MessageState) String() string {
	return string(s)
}

type Message struct {
	id        MessageID
	mType     MessageType
	sender    kernel.UserID
	recipient kernel.ID
	state     MessageState
	content   string
	sentAt    time.Time
}

func NewMessage(
	id MessageID,
	mType MessageType,
	recipient kernel.ID,
	state MessageState,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
) *Message {
	return &Message{
		id:        id,
		mType:     mType,
		recipient: recipient,
		state:     state,
		sender:    sender,
		content:   content,
		sentAt:    sentAt,
	}
}

func (m *Message) BeDelivered() {
	m.state = MessageStateDelivered
}

func (m *Message) BeRead() {
	m.state = MessageStateRead
}

func (m *Message) ID() MessageID {
	return m.id
}

func (m *Message) Type() MessageType {
	return m.mType
}

func (m *Message) Recipient() kernel.ID {
	return m.recipient
}

func (m *Message) State() MessageState {
	return m.state
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
