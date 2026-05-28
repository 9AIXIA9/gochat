package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicRoomMessageCreated event.Topic = "room_message.created"

var _ event.SpecificEvent = (*RoomMessageCreatedEvent)(nil)

type RoomMessageCreatedEvent struct {
	senderID     kernel.UserID
	recipientIDs []kernel.UserID
	roomID       kernel.RoomID
	content      string
	sentAt       time.Time
	*event.StandardEvent
}

func NewRoomMessageCreatedEvent(
	messageID kernel.MessageID,
	senderID kernel.UserID,
	recipientIDs []kernel.UserID,
	roomID kernel.RoomID,
	content string,
	sentAt time.Time,
	generator event.IDGenerator,
) (*RoomMessageCreatedEvent, error) {
	e := &RoomMessageCreatedEvent{
		senderID:     senderID,
		recipientIDs: recipientIDs,
		roomID:       roomID,
		content:      content,
		sentAt:       sentAt,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(messageID), TopicRoomMessageCreated, payload, generator)
	return e, nil
}

func (e *RoomMessageCreatedEvent) SenderID() kernel.UserID {
	return e.senderID
}

func (e *RoomMessageCreatedEvent) RecipientIDs() []kernel.UserID {
	return e.recipientIDs
}

func (e *RoomMessageCreatedEvent) RoomID() kernel.RoomID {
	return e.roomID
}

func (e *RoomMessageCreatedEvent) Content() string {
	return e.content
}

func (e *RoomMessageCreatedEvent) SentAt() time.Time {
	return e.sentAt
}

func (e *RoomMessageCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		SenderID     kernel.UserID
		RecipientIDs []kernel.UserID
		RoomID       kernel.RoomID
		Content      string
		SentAt       time.Time
	}
	return json.Marshal(&Alias{
		SenderID:     e.senderID,
		RecipientIDs: e.recipientIDs,
		RoomID:       e.roomID,
		Content:      e.content,
		SentAt:       e.sentAt,
	})
}

func (e *RoomMessageCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		SenderID     kernel.UserID
		RecipientIDs []kernel.UserID
		RoomID       kernel.RoomID
		Content      string
		SentAt       time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.senderID = tmp.SenderID
	e.recipientIDs = tmp.RecipientIDs
	e.roomID = tmp.RoomID
	e.content = tmp.Content
	e.sentAt = tmp.SentAt
	return nil
}
