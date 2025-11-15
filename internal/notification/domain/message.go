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
	id      MessageID
	sender  kernel.UserID
	state   MessageState
	content string
	sentAt  time.Time
}

func NewMessage(
	id MessageID,
	sender kernel.UserID,
	state MessageState,
	content string,
	sentAt time.Time,
) *Message {
	return &Message{
		id:      id,
		sender:  sender,
		state:   state,
		content: content,
		sentAt:  sentAt,
	}
}

func (m *Message) ID() MessageID {
	return m.id
}

func (m *Message) Sender() kernel.UserID {
	return m.sender
}

func (m *Message) State() MessageState {
	return m.state
}

func (m *Message) Content() string {
	return m.content
}

func (m *Message) SentAt() time.Time {
	return m.sentAt
}
