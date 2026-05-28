package domain

import (
	"context"
	"encoding/json"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

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
	ctx context.Context,
	recipientID kernel.UserID,
	senderID kernel.UserID,
	content string,
	eventIDGenerator event.IDGenerator,
	messageIDGenerator kernel.MessageIDGenerator,
	exister FriendshipExisterByUserID,
) (*PrivateMessage, error) {
	if len(content) == 0 {
		return nil, ErrEmptyMessageContent
	}

	if senderID != recipientID {
		exist, err := exister.ExistByUserID(ctx, senderID, recipientID)
		if err != nil {
			return nil, err
		}

		if !exist {
			return nil, ErrNotFriends
		}
	}

	message := &PrivateMessage{
		id:           messageIDGenerator.Generate(),
		senderID:     senderID,
		recipientID:  recipientID,
		content:      content,
		sentAt:       time.Now().UTC(),
		eventManager: event.NewEventManager(),
	}

	rawPayload, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	ev, err := contract.NewNotificationCreatedEvent(
		message.id,
		message.recipientID,
		rawPayload,
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

func (m *PrivateMessage) Marshal() ([]byte, error) {
	type Alias struct {
		ID          kernel.MessageID
		SenderID    kernel.UserID
		RecipientID kernel.UserID
		Content     string
		SentAt      time.Time
	}

	return json.Marshal(&Alias{
		ID:          m.id,
		SenderID:    m.senderID,
		RecipientID: m.recipientID,
		Content:     m.content,
		SentAt:      m.sentAt,
	})
}

func (m *PrivateMessage) GetEvents() []event.Event {
	return m.eventManager.GetEvents()
}
