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
	content string
	sentAt  time.Time
	states  []*RecipientMessageState
}

func NewMessage(
	id MessageID,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
	states []*RecipientMessageState,
) *Message {
	return &Message{
		id:      id,
		sender:  sender,
		content: content,
		sentAt:  sentAt,
		states:  states,
	}
}

func (m *Message) ID() MessageID {
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

func (m *Message) States() []*RecipientMessageState {
	return m.states
}

func (m *Message) Recipients() []kernel.UserID {
	recipients := make([]kernel.UserID, 0, len(m.states))
	for _, state := range m.states {
		recipients = append(recipients, state.recipient)
	}
	return recipients
}
