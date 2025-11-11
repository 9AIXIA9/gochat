package domain

import (
	"encoding/json"
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

func (m *Message) Marshal() ([]byte, error) {
	type Alias struct {
		ID      MessageID
		Sender  kernel.UserID
		State   MessageState
		Content string
		SentAt  time.Time
	}
	return json.Marshal(Alias{
		ID:      m.id,
		Sender:  m.sender,
		State:   m.state,
		Content: m.content,
		SentAt:  m.sentAt,
	})
}

func (m *Message) Unmarshal(data []byte) error {
	type Alias struct {
		ID      MessageID
		Sender  kernel.UserID
		State   MessageState
		Content string
		SentAt  time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	m.id = tmp.ID
	m.sender = tmp.Sender
	m.state = tmp.State
	m.content = tmp.Content
	m.sentAt = tmp.SentAt
	return nil
}
