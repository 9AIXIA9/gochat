package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicPrivateMessageNotificationRequested event.Topic = "notification.private_message_notification.requested"

var _ event.SpecificEvent = (*PrivateMessageNotificationRequestedEvent)(nil)

type PrivateMessageNotificationRequestedEvent struct {
	senderID    kernel.UserID
	recipientID kernel.UserID
	content     string
	sentAt      time.Time
	*event.StandardEvent
}

func ToPrivateMessageNotificationRequestedEvent(ev event.Event) (*PrivateMessageNotificationRequestedEvent, error) {
	e := &PrivateMessageNotificationRequestedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewPrivateMessageNotificationRequestedEvent(
	id event.ID,
	messageID kernel.MessageID,
	recipientID kernel.UserID,
	senderID kernel.UserID,
	content string,
	sentAt time.Time,
) (*PrivateMessageNotificationRequestedEvent, error) {
	e := &PrivateMessageNotificationRequestedEvent{
		senderID:    senderID,
		recipientID: recipientID,
		content:     content,
		sentAt:      sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(messageID), time.Now().UTC(), TopicPrivateMessageNotificationRequested, payload)
	return e, nil
}

func (e *PrivateMessageNotificationRequestedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID   kernel.MessageID
		SenderID    kernel.UserID
		RecipientID kernel.UserID
		Content     string
		SentAt      time.Time
	}
	return json.Marshal(Alias{
		SenderID:    e.senderID,
		RecipientID: e.recipientID,
		Content:     e.content,
		SentAt:      e.sentAt,
	})
}

func (e *PrivateMessageNotificationRequestedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID   kernel.MessageID
		SenderID    kernel.UserID
		RecipientID kernel.UserID
		Content     string
		SentAt      time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.senderID = tmp.SenderID
	e.content = tmp.Content
	e.sentAt = tmp.SentAt
	e.recipientID = tmp.RecipientID
	return nil
}

func (e *PrivateMessageNotificationRequestedEvent) SenderID() kernel.UserID {
	return e.senderID
}

func (e *PrivateMessageNotificationRequestedEvent) RecipientID() kernel.UserID {
	return e.recipientID
}

func (e *PrivateMessageNotificationRequestedEvent) Content() string {
	return e.content
}

func (e *PrivateMessageNotificationRequestedEvent) SentAt() time.Time {
	return e.sentAt
}
