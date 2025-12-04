package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicRoomCreated event.Topic = "roomship.room.created"

var _ event.SpecificEvent = (*RoomCreatedEvent)(nil)

type RoomCreatedEvent struct {
	*event.StandardEvent
}

func ToRoomCreatedEvent(ev event.Event) (*RoomCreatedEvent, error) {
	if ev.Topic() != TopicRoomCreated {
		return nil, myErrors.ErrWrongEventType
	}
	e := &RoomCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomCreatedEvent(
	roomID kernel.RoomID,
	generator event.IDGenerator,
) (*RoomCreatedEvent, error) {
	e := &RoomCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(roomID), TopicRoomCreated, payload, generator)
	return e, nil
}

func (e *RoomCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *RoomCreatedEvent) Unmarshal(_ []byte) error {
	return nil
}
