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
	messageID MessageID
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
) (*RoomMessageReceivedEvent, error) {
	e := &RoomMessageReceivedEvent{
		messageID: messageID,
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
		MessageID MessageID
	}
	return json.Marshal(Alias{
		MessageID: e.messageID,
	})
}

func (e *RoomMessageReceivedEvent) Unmarshal(data []byte) error {
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

func (e *RoomMessageReceivedEvent) MessageID() MessageID {
	return e.messageID
}
