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
	state   MessageState
	sender  kernel.UserID
	content string
	sentAt  time.Time
}

func NewMessage(id MessageID, state MessageState, sender kernel.UserID, content string, sentAt time.Time) *Message {
	return &Message{id: id, state: state, sender: sender, content: content, sentAt: sentAt}
}

func (m *Message) ID() MessageID {
	return m.id
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
