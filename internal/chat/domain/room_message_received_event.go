package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicRoomMessageReceived event.Topic = "chat.room_message.received"

var _ event.SpecificEvent = (*RoomMessageReceivedEvent)(nil)

type RoomMessageReceivedEvent struct {
	messageID  MessageID
	roomID     RoomID
	sender     kernel.UserID
	recipients []kernel.UserID
	content    string
	sentAt     time.Time
	*event.StandardEvent
}

func ToRoomMessageReceivedEvent(ev event.Event) (*RoomMessageReceivedEvent, error) {
	e := &RoomMessageReceivedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomMessageReceivedEvent(
	id event.ID,
	messageID MessageID,
	roomID RoomID,
	sender kernel.UserID,
	recipients []kernel.UserID,
	content string,
	sentAt time.Time,
) (*RoomMessageReceivedEvent, error) {
	e := &RoomMessageReceivedEvent{
		messageID:  messageID,
		roomID:     roomID,
		sender:     sender,
		recipients: recipients,
		content:    content,
		sentAt:     sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(roomID), time.Now().UTC(), TopicRoomMessageReceived, payload)
	return e, nil
}

func (e *RoomMessageReceivedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID  MessageID
		RoomID     RoomID
		Sender     kernel.UserID
		Recipients []kernel.UserID
		Content    string
		SentAt     time.Time
	}
	return json.Marshal(Alias{
		MessageID:  e.messageID,
		RoomID:     e.roomID,
		Sender:     e.sender,
		Recipients: e.recipients,
		Content:    e.content,
		SentAt:     e.sentAt,
	})
}

func (e *RoomMessageReceivedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID  MessageID
		RoomID     RoomID
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
	e.roomID = tmp.RoomID
	e.sender = tmp.Sender
	e.recipients = tmp.Recipients
	e.content = tmp.Content
	e.sentAt = tmp.SentAt
	return nil
}

func (e *RoomMessageReceivedEvent) MessageID() MessageID {
	return e.messageID
}

func (e *RoomMessageReceivedEvent) Sender() kernel.UserID {
	return e.sender
}

func (e *RoomMessageReceivedEvent) RoomID() RoomID {
	return e.roomID
}

func (e *RoomMessageReceivedEvent) Content() string {
	return e.content
}

func (e *RoomMessageReceivedEvent) SentAt() time.Time {
	return e.sentAt
}
