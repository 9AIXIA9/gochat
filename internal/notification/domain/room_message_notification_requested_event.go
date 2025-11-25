package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicRoomMessageNotificationRequested event.Topic = "notification.room_message_notification.requested"

var _ event.SpecificEvent = (*RoomMessageNotificationRequestedEvent)(nil)

type RoomMessageNotificationRequestedEvent struct {
	messageID    kernel.MessageID
	senderID     kernel.UserID
	roomID       kernel.RoomID
	recipientIDs []kernel.UserID
	content      string
	sentAt       time.Time
	*event.StandardEvent
}

func ToRoomMessageNotificationRequestedEvent(ev event.Event) (*RoomMessageNotificationRequestedEvent, error) {
	if ev.Topic() != TopicRoomMessageNotificationRequested {
		return nil, myErrors.ErrWrongEventType
	}
	e := &RoomMessageNotificationRequestedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomMessageNotificationRequestedEvent(
	messageID kernel.MessageID,
	senderID kernel.UserID,
	roomID kernel.RoomID,
	recipientIDs []kernel.UserID,
	content string,
	sentAt time.Time,
	generator event.IDGenerator,
) (*RoomMessageNotificationRequestedEvent, error) {
	e := &RoomMessageNotificationRequestedEvent{
		messageID:    messageID,
		senderID:     senderID,
		roomID:       roomID,
		recipientIDs: recipientIDs,
		content:      content,
		sentAt:       sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(messageID), TopicRoomMessageNotificationRequested, payload, generator)
	return e, nil
}

func (e *RoomMessageNotificationRequestedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID    kernel.MessageID
		SenderID     kernel.UserID
		RoomID       kernel.RoomID
		RecipientIDs []kernel.UserID
		Content      string
		SentAt       time.Time
	}
	return json.Marshal(Alias{
		MessageID:    e.messageID,
		SenderID:     e.senderID,
		RoomID:       e.roomID,
		RecipientIDs: e.recipientIDs,
		Content:      e.content,
		SentAt:       e.sentAt,
	})
}

func (e *RoomMessageNotificationRequestedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID    kernel.MessageID
		SenderID     kernel.UserID
		RoomID       kernel.RoomID
		RecipientIDs []kernel.UserID
		Content      string
		SentAt       time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.messageID = tmp.MessageID
	e.senderID = tmp.SenderID
	e.roomID = tmp.RoomID
	e.content = tmp.Content
	e.sentAt = tmp.SentAt
	e.recipientIDs = tmp.RecipientIDs
	return nil
}

func (e *RoomMessageNotificationRequestedEvent) SenderID() kernel.UserID {
	return e.senderID
}

func (e *RoomMessageNotificationRequestedEvent) RoomID() kernel.RoomID {
	return e.roomID
}

func (e *RoomMessageNotificationRequestedEvent) RecipientIDs() []kernel.UserID {
	return e.recipientIDs
}

func (e *RoomMessageNotificationRequestedEvent) Content() string {
	return e.content
}

func (e *RoomMessageNotificationRequestedEvent) SentAt() time.Time {
	return e.sentAt
}
