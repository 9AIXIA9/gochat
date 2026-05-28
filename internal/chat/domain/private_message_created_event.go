package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicPrivateMessageCreated event.Topic = "private_message.created"

var _ event.SpecificEvent = (*PrivateMessageCreatedEvent)(nil)

type PrivateMessageCreatedEvent struct {
	*event.StandardEvent
	senderID    kernel.UserID
	recipientID kernel.UserID
	content     string
	sentAt      time.Time
}

func NewPrivateMessageCreatedEvent(
	messageID kernel.MessageID,
	senderID kernel.UserID,
	recipientID kernel.UserID,
	content string,
	sentAt time.Time,
	generator event.IDGenerator,
) (*PrivateMessageCreatedEvent, error) {
	e := &PrivateMessageCreatedEvent{
		senderID:    senderID,
		recipientID: recipientID,
		content:     content,
		sentAt:      sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(messageID), TopicPrivateMessageCreated, payload, generator)
	return e, nil
}

func (e *PrivateMessageCreatedEvent) SenderID() kernel.UserID {
	return e.senderID
}

func (e *PrivateMessageCreatedEvent) RecipientID() kernel.UserID {
	return e.recipientID
}

func (e *PrivateMessageCreatedEvent) Content() string {
	return e.content
}

func (e *PrivateMessageCreatedEvent) SentAt() time.Time {
	return e.sentAt
}

func (e *PrivateMessageCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		SenderID    kernel.UserID
		RecipientID kernel.UserID
		Content     string
		SentAt      time.Time
	}
	return json.Marshal(&Alias{
		SenderID:    e.senderID,
		RecipientID: e.recipientID,
		Content:     e.content,
		SentAt:      e.sentAt,
	})
}

func (e *PrivateMessageCreatedEvent) Unmarshal([]byte) error {
	type Alias struct {
		SenderID    kernel.UserID
		RecipientID kernel.UserID
		Content     string
		SentAt      time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(e.Payload(), &tmp); err != nil {
		return err
	}
	e.senderID = tmp.SenderID
	e.recipientID = tmp.RecipientID
	e.content = tmp.Content
	e.sentAt = tmp.SentAt
	return nil
}
