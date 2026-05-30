package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
)

const TopicRoomCreated event.Topic = "room.created"

var _ event.SpecificEvent = (*RoomCreatedEvent)(nil)

type RoomCreatedEvent struct {
	*event.StandardEvent
}

func ToRoomCreatedEvent(ev event.Event) (*RoomCreatedEvent, error) {
	if ev.Topic() != TopicRoomCreated {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &RoomCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *RoomCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *RoomCreatedEvent) Unmarshal([]byte) error {
	return nil
}
