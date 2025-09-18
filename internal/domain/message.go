package domain

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type BaseNumber int64

// Message 泛型消息结构
type Message struct {
	id      string
	from    UserNumber
	to      BaseNumber
	content string
	sentAt  time.Time
	sent    bool
	//read  bool
}

func NewMessage(id string, from UserNumber, to BaseNumber, content string, sendAt time.Time, sent bool) *Message {
	return &Message{
		id:      id,
		from:    from,
		to:      to,
		content: content,
		sentAt:  sendAt,
		sent:    sent,
	}
}

func CreateMessage(from UserNumber, to BaseNumber, content string, sendAt time.Time, sent bool) *Message {
	return NewMessage(uuid.NewString(), from, to, content, sendAt, sent)
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

func (m *Message) ID() string {
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

func (m *Message) IsSent() bool {
	return m.sent
}

func (m *Message) Sent() {
	m.sent = true
}
