package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicRoomMessageCreated event.Topic = "chat.private_message.created"

var _ event.SpecificEvent = (*RoomMessageCreatedEvent)(nil)

type RoomMessageCreatedEvent struct {
	messageID  MessageID
	sender     kernel.UserID
	recipients []kernel.UserID
	content    string
	sentAt     time.Time
	*event.StandardEvent
}

func ToRoomMessageCreatedEvent(ev event.Event) (*RoomMessageCreatedEvent, error) {
	e := &RoomMessageCreatedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomMessageCreatedEvent(
	id event.ID,
	messageID MessageID,
	sender kernel.UserID,
	recipients []kernel.UserID,
	content string,
	sentAt time.Time,
) (*RoomMessageCreatedEvent, error) {
	e := &RoomMessageCreatedEvent{
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

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(sender), time.Now().UTC(), TopicRoomMessageCreated, payload)
	return e, nil
}

func (e *RoomMessageCreatedEvent) Marshal() ([]byte, error) {
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

func (e *RoomMessageCreatedEvent) Unmarshal(data []byte) error {
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

func (e *RoomMessageCreatedEvent) MessageID() MessageID {
	return e.messageID
}

func (e *RoomMessageCreatedEvent) Sender() kernel.UserID {
	return e.sender
}

func (e *RoomMessageCreatedEvent) Recipients() []kernel.UserID {
	return e.recipients
}

func (e *RoomMessageCreatedEvent) Content() string {
	return e.content
}

func (e *RoomMessageCreatedEvent) SentAt() time.Time {
	return e.sentAt
}
