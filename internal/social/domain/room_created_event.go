package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicRoomCreated event.Topic = "social.room.created"

var _ event.SpecificEvent = (*RoomCreatedEvent)(nil)

type RoomCreatedEvent struct {
	roomID RoomID
	*event.StandardEvent
}

func ToRoomCreatedEvent(ev event.Event) (*RoomCreatedEvent, error) {
	e := &RoomCreatedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomCreatedEvent(id event.ID, roomID RoomID) (*RoomCreatedEvent, error) {
	e := &RoomCreatedEvent{
		roomID: roomID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(roomID), time.Now().UTC(), TopicRoomCreated, payload)
	return e, nil
}

func (e *RoomCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		RoomID RoomID
	}
	return json.Marshal(Alias{
		RoomID: e.roomID,
	})
}

func (e *RoomCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		RoomID RoomID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.roomID = tmp.RoomID
	return nil
}
