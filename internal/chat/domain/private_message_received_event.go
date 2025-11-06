package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicPrivateMessageReceived event.Topic = "chat.private_message.received"

var _ event.SpecificEvent = (*PrivateMessageReceivedEvent)(nil)

type PrivateMessageReceivedEvent struct {
	messageID MessageID
	recipient kernel.UserID
	sender    kernel.UserID
	content   string
	sentAt    time.Time
	*event.StandardEvent
}

func ToPrivateMessageReceivedEvent(ev event.Event) (*PrivateMessageReceivedEvent, error) {
	e := &PrivateMessageReceivedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewPrivateMessageReceivedEvent(
	id event.ID,
	messageID MessageID,
	recipient kernel.UserID,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
) (*PrivateMessageReceivedEvent, error) {
	e := &PrivateMessageReceivedEvent{
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

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(recipient), time.Now().UTC(), TopicPrivateMessageReceived, payload)
	return e, nil
}

func (e *PrivateMessageReceivedEvent) Marshal() ([]byte, error) {
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

func (e *PrivateMessageReceivedEvent) Unmarshal(data []byte) error {
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

func (e *PrivateMessageReceivedEvent) MessageID() MessageID {
	return e.messageID
}

func (e *PrivateMessageReceivedEvent) Sender() kernel.UserID {
	return e.sender
}

func (e *PrivateMessageReceivedEvent) Recipient() kernel.UserID {
	return e.recipient
}

func (e *PrivateMessageReceivedEvent) Content() string {
	return e.content
}

func (e *PrivateMessageReceivedEvent) SentAt() time.Time {
	return e.sentAt
}
