package domain

import (
	"encoding/json"
	"time"
)

type BaseNumber int64

// Message 泛型消息结构
type Message struct {
	from    UserNumber
	to      BaseNumber
	content string
	sentAt  time.Time
	//todo
	//isRead  bool
	//isSent  bool
}

func NewMessage(from UserNumber, to BaseNumber, content string, sendAt time.Time) *Message {
	return &Message{
		from:    from,
		to:      to,
		content: content,
		sentAt:  sendAt,
	}
}

func (m *Message) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"from":    m.from,
		"to":      m.to,
		"content": m.content,
		"sent_at": m.sentAt.Format(time.RFC3339),
	}
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToJSON())
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
