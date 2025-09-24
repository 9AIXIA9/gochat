package domain

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type BaseNumber int64
type MessageID string

type Message struct {
	id      MessageID
	from    UserNumber
	to      BaseNumber
	content string
	sentAt  time.Time
}

func NewMessage(id MessageID, from UserNumber, to BaseNumber, content string, sendAt time.Time) *Message {
	return &Message{
		id:      id,
		from:    from,
		to:      to,
		content: content,
		sentAt:  sendAt,
	}
}

func CreateMessage(from UserNumber, to BaseNumber, content string, sendAt time.Time) *Message {
	return NewMessage(MessageID(uuid.NewString()), from, to, content, sendAt)
}

func (m *Message) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":      m.id,
		"from":    m.from,
		"to":      m.to,
		"content": m.content,
		"sent_at": m.sentAt.Format(time.RFC3339),
	}
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToJSON())
}

func (m *Message) ID() MessageID {
	return m.id
}

func (m *Message) From() UserNumber {
	return m.from
}

func (m *Message) To() BaseNumber {
	return m.to
}

func (m *Message) Content() string {
	return m.content
}

func (m *Message) SendAt() time.Time {
	return m.sentAt
}
