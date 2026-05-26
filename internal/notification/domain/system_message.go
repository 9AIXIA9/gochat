package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type SystemMessage struct {
	id          kernel.MessageID
	recipientID kernel.UserID
	state       MessageState
	content     string
	sentAt      time.Time
}

func LoadSystemMessage(
	id kernel.MessageID,
	recipientID kernel.UserID,
	state MessageState,
	content string,
	sentAt time.Time,
) *SystemMessage {
	return &SystemMessage{
		id:          id,
		recipientID: recipientID,
		state:       state,
		content:     content,
		sentAt:      sentAt,
	}
}

func CreateSystemMessage(
	recipientID kernel.UserID,
	content string,
	idGenerator kernel.MessageIDGenerator,
) (*SystemMessage, error) {
	if len(content) == 0 {
		return nil, ErrEmptyContent
	}

	return &SystemMessage{
		id:          idGenerator.Generate(),
		recipientID: recipientID,
		state:       MessageStateUndelivered,
		content:     content,
		sentAt:      time.Now().UTC(),
	}, nil
}

func (m *SystemMessage) ID() kernel.MessageID {
	return m.id
}

func (m *SystemMessage) RecipientID() kernel.UserID {
	return m.recipientID
}

func (m *SystemMessage) State() MessageState {
	return m.state
}

func (m *SystemMessage) Content() string {
	return m.content
}

func (m *SystemMessage) SentAt() time.Time {
	return m.sentAt
}
