package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

//TODO 重设计 多聚合根关系由服务层进行协调 领域层只负责单一聚合根的业务逻辑

type PrivateMessage struct {
	id          kernel.MessageID
	senderID    kernel.UserID
	recipientID kernel.UserID
	content     string
	sentAt      time.Time

	eventManager *event.Manager
}

func LoadPrivateMessage(
	id kernel.MessageID,
	senderID kernel.UserID,
	recipientID kernel.UserID,
	content string,
	sentAt time.Time,
) *PrivateMessage {
	return &PrivateMessage{
		id:           id,
		senderID:     senderID,
		recipientID:  recipientID,
		content:      content,
		sentAt:       sentAt,
		eventManager: event.NewEventManager(),
	}
}

func CreatePrivateMessage(
	recipient *User,
	senderID kernel.UserID,
	content string,
	messageIDGenerator MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
) (*PrivateMessage, error) {
	message := &PrivateMessage{
		id:           messageIDGenerator.Generate(),
		senderID:     senderID,
		recipientID:  recipient.id,
		content:      content,
		sentAt:       time.Now().UTC(),
		eventManager: event.NewEventManager(),
	}

	ev, err := NewPrivateMessageCreatedEvent(
		message.id,
		eventIDGenerator,
	)
	if err != nil {
		return nil, err
	}

	message.eventManager.RecordEvent(ev)

	return message, nil
}

func (m *PrivateMessage) ID() kernel.MessageID {
	return m.id
}

func (m *PrivateMessage) SenderID() kernel.UserID {
	return m.senderID
}

func (m *PrivateMessage) RecipientID() kernel.UserID {
	return m.recipientID
}

func (m *PrivateMessage) Content() string {
	return m.content
}

func (m *PrivateMessage) SentAt() time.Time {
	return m.sentAt
}

func (m *PrivateMessage) GetEvents() []event.Event {
	return m.eventManager.GetEvents()
}
