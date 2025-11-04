package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicMessageSent event.Topic = "chat.message.sent"

var _ event.SpecificEvent = (*MessageSentEvent)(nil)

type MessageSentEvent struct {
	sender    kernel.UserID
	recipient kernel.ID
	content   string
	mType     MessageType
	sentAt    time.Time
	*event.StandardEvent
}

func ToMessageSentEvent(ev event.Event) (*MessageSentEvent, error) {
	e := &MessageSentEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewMessageSentEvent(
	id event.ID,
	sender kernel.UserID,
	recipient kernel.ID,
	content string,
	mType MessageType,
	sentAt time.Time,
) (*MessageSentEvent, error) {
	e := &MessageSentEvent{
		sender:    sender,
		recipient: recipient,
		content:   content,
		mType:     mType,
		sentAt:    sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(sender), time.Now().UTC(), TopicMessageSent, payload)
	return e, nil
}

func (e *MessageSentEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Sender    kernel.UserID
		Recipient kernel.ID
		Content   string
		Type      MessageType
		SentAt    time.Time
	}
	return json.Marshal(Alias{
		Sender:    e.sender,
		Recipient: e.recipient,
		Content:   e.content,
		Type:      e.mType,
		SentAt:    e.sentAt,
	})
}

func (e *MessageSentEvent) Unmarshal(data []byte) error {
	type Alias struct {
		Sender    kernel.UserID
		Recipient kernel.ID
		Content   string
		Type      MessageType
		SentAt    time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.sender = tmp.Sender
	e.recipient = tmp.Recipient
	e.content = tmp.Content
	e.mType = tmp.Type
	e.sentAt = tmp.SentAt
	return nil
}

func (e *MessageSentEvent) Sender() kernel.UserID {
	return e.sender
}

func (e *MessageSentEvent) Recipient() kernel.ID {
	return e.recipient
}

func (e *MessageSentEvent) Content() string {
	return e.content
}

func (e *MessageSentEvent) Type() MessageType {
	return e.mType
}

func (e *MessageSentEvent) SentAt() time.Time {
	return e.sentAt
}
