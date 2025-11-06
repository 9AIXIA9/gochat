package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicMessageReceived event.Topic = "chat.message.received"

var _ event.SpecificEvent = (*MessageReceivedEvent)(nil)

type MessageReceivedEvent struct {
	messageID MessageID
	recipient kernel.UserID
	sender    kernel.UserID
	content   string
	sentAt    time.Time
	*event.StandardEvent
}

func ToMessageReceivedEvent(ev event.Event) (*MessageReceivedEvent, error) {
	e := &MessageReceivedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewMessageReceivedEvent(
	id event.ID,
	messageID MessageID,
	recipient kernel.UserID,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
) (*MessageReceivedEvent, error) {
	e := &MessageReceivedEvent{
		messageID: messageID,
		recipient: recipient,
		sender:    sender,
		content:   content,
		sentAt:    sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(recipient), time.Now().UTC(), TopicMessageReceived, payload)
	return e, nil
}

func (e *MessageReceivedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID MessageID
		Recipient kernel.UserID
		Sender    kernel.UserID
		Content   string
		SentAt    time.Time
	}
	return json.Marshal(Alias{
		MessageID: e.messageID,
		Recipient: e.recipient,
		Sender:    e.sender,
		Content:   e.content,
		SentAt:    e.sentAt,
	})
}

func (e *MessageReceivedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID MessageID
		Recipient kernel.UserID
		Sender    kernel.UserID
		Content   string
		SentAt    time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.messageID = tmp.MessageID
	e.recipient = tmp.Recipient
	e.sender = tmp.Sender
	e.content = tmp.Content
	e.sentAt = tmp.SentAt
	return nil
}

func (e *MessageReceivedEvent) MessageID() MessageID {
	return e.messageID
}

func (e *MessageReceivedEvent) Sender() kernel.UserID {
	return e.sender
}

func (e *MessageReceivedEvent) Recipient() kernel.UserID {
	return e.recipient
}

func (e *MessageReceivedEvent) Content() string {
	return e.content
}

func (e *MessageReceivedEvent) SentAt() time.Time {
	return e.sentAt
}
