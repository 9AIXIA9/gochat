package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicRoomCreated event.Topic = "social.room.created"

var _ event.SpecificEvent = (*RoomCreatedEvent)(nil)

type RoomCreatedEvent struct {
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
	e := &RoomCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(roomID), time.Now().UTC(), TopicRoomCreated, payload)
	return e, nil
}

func (e *RoomCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *RoomCreatedEvent) Unmarshal(_ []byte) error {
	return nil
}
