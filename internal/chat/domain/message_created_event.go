package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicMessageCreated event.Topic = "chat.message.created"

var _ event.SpecificEvent = (*MessageCreatedEvent)(nil)

type MessageCreatedEvent struct {
	messageID  MessageID
	sender     kernel.UserID
	recipients []kernel.UserID
	content    string
	sentAt     time.Time
	*event.StandardEvent
}

func ToMessageCreatedEvent(ev event.Event) (*MessageCreatedEvent, error) {
	e := &MessageCreatedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewMessageCreatedEvent(
	id event.ID,
	messageID MessageID,
	sender kernel.UserID,
	recipients []kernel.UserID,
	content string,
	sentAt time.Time,
) (*MessageCreatedEvent, error) {
	e := &MessageCreatedEvent{
		messageID:  messageID,
		sender:     sender,
		recipients: recipients,
		content:    content,
		sentAt:     sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(sender), time.Now().UTC(), TopicMessageCreated, payload)
	return e, nil
}

func (e *MessageCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID  MessageID
		Sender     kernel.UserID
		Recipients []kernel.UserID
		Content    string
		SentAt     time.Time
	}
	return json.Marshal(Alias{
		MessageID:  e.messageID,
		Sender:     e.sender,
		Recipients: e.recipients,
		Content:    e.content,
		SentAt:     e.sentAt,
	})
}

func (e *MessageCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID  MessageID
		Sender     kernel.UserID
		Recipients []kernel.UserID
		Content    string
		SentAt     time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.messageID = tmp.MessageID
	e.sender = tmp.Sender
	e.recipients = tmp.Recipients
	e.content = tmp.Content
	e.sentAt = tmp.SentAt
	return nil
}

func (e *MessageCreatedEvent) MessageID() MessageID {
	return e.messageID
}

func (e *MessageCreatedEvent) Sender() kernel.UserID {
	return e.sender
}

func (e *MessageCreatedEvent) Recipients() []kernel.UserID {
	return e.recipients
}

func (e *MessageCreatedEvent) Content() string {
	return e.content
}

func (e *MessageCreatedEvent) SentAt() time.Time {
	return e.sentAt
}
