package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicMessageNotificationRequested event.Topic = "notification.message_notification.requested"

var _ event.SpecificEvent = (*MessageNotificationRequestedEvent)(nil)

type MessageNotificationRequestedEvent struct {
	messageID kernel.MessageID
	sender    kernel.UserID
	content   string
	sentAt    time.Time
	*event.StandardEvent
}

func ToMessageNotificationRequestedEvent(ev event.Event) (*MessageNotificationRequestedEvent, error) {
	e := &MessageNotificationRequestedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewMessageNotificationRequestedEvent(
	id event.ID,
	messageID kernel.MessageID,
	recipient kernel.UserID,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
) (*MessageNotificationRequestedEvent, error) {
	e := &MessageNotificationRequestedEvent{
		messageID: messageID,
		sender:    sender,
		content:   content,
		sentAt:    sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(recipient), time.Now().UTC(), TopicMessageNotificationRequested, payload)
	return e, nil
}

func (e *MessageNotificationRequestedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID kernel.MessageID
		Sender    kernel.UserID
		Content   string
		SentAt    time.Time
	}
	return json.Marshal(Alias{
		MessageID: e.messageID,
		Sender:    e.sender,
		Content:   e.content,
		SentAt:    e.sentAt,
	})
}

func (e *MessageNotificationRequestedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID kernel.MessageID
		Sender    kernel.UserID
		Content   string
		SentAt    time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.messageID = tmp.MessageID
	e.sender = tmp.Sender
	e.content = tmp.Content
	e.sentAt = tmp.SentAt
	return nil
}

func (e *MessageNotificationRequestedEvent) MessageID() kernel.MessageID {
	return e.messageID
}

func (e *MessageNotificationRequestedEvent) Sender() kernel.UserID {
	return e.sender
}

func (e *MessageNotificationRequestedEvent) Content() string {
	return e.content
}

func (e *MessageNotificationRequestedEvent) SentAt() time.Time {
	return e.sentAt
}
