package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicRoomMessageCreated event.Topic = "notification.room_message.created"

var _ event.SpecificEvent = (*RoomMessageCreatedEvent)(nil)

type RoomMessageCreatedEvent struct {
	messageID MessageID
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
	roomID RoomID,
) (*RoomMessageCreatedEvent, error) {
	e := &RoomMessageCreatedEvent{
		messageID: messageID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(roomID), time.Now().UTC(), TopicRoomMessageCreated, payload)
	return e, nil
}

func (e *RoomMessageCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		MessageID MessageID
	}
	return json.Marshal(Alias{
		MessageID: e.messageID,
	})
}

func (e *RoomMessageCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		MessageID MessageID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.messageID = tmp.MessageID
	return nil
}

func (e *RoomMessageCreatedEvent) MessageID() MessageID {
	return e.messageID
}
