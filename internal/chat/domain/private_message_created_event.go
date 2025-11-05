package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicPrivateMessageCreated event.Topic = "chat.private_message.created"

var _ event.SpecificEvent = (*PrivateMessageCreatedEvent)(nil)

type PrivateMessageCreatedEvent struct {
	messageID MessageID
	sender    kernel.UserID
	recipient kernel.UserID
	content   string
	sentAt    time.Time
	*event.StandardEvent
}

func ToPrivateMessageCreatedEvent(ev event.Event) (*PrivateMessageCreatedEvent, error) {
	e := &PrivateMessageCreatedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewPrivateMessageCreatedEvent(
	id event.ID,
	messageID MessageID,
	sender kernel.UserID,
	recipient kernel.UserID,
	content string,
	sentAt time.Time,
) (*PrivateMessageCreatedEvent, error) {
	e := &PrivateMessageCreatedEvent{
		messageID: messageID,
		sender:    sender,
		recipient: recipient,
		content:   content,
		sentAt:    sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(sender), time.Now().UTC(), TopicPrivateMessageCreated, payload)
	return e, nil
}

func (e *PrivateMessageCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID MessageID
		Sender    kernel.UserID
		Recipient kernel.UserID
		Content   string
		SentAt    time.Time
	}
	return json.Marshal(Alias{
		MessageID: e.messageID,
		Sender:    e.sender,
		Recipient: e.recipient,
		Content:   e.content,
		SentAt:    e.sentAt,
	})
}

func (e *PrivateMessageCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID MessageID
		Sender    kernel.UserID
		Recipient kernel.UserID
		Content   string
		SentAt    time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.messageID = tmp.MessageID
	e.sender = tmp.Sender
	e.recipient = tmp.Recipient
	e.content = tmp.Content
	e.sentAt = tmp.SentAt
	return nil
}

func (e *PrivateMessageCreatedEvent) MessageID() MessageID {
	return e.messageID
}

func (e *PrivateMessageCreatedEvent) Sender() kernel.UserID {
	return e.sender
}

func (e *PrivateMessageCreatedEvent) Recipient() kernel.UserID {
	return e.recipient
}

func (e *PrivateMessageCreatedEvent) Content() string {
	return e.content
}

func (e *PrivateMessageCreatedEvent) SentAt() time.Time {
	return e.sentAt
}
