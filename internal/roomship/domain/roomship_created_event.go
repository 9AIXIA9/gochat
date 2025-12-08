package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicRoomshipCreated event.Topic = "roomship.roomship.created"

var _ event.SpecificEvent = (*RoomshipCreatedEvent)(nil)

type RoomshipCreatedEvent struct {
	*event.StandardEvent
}

func ToRoomshipCreatedEvent(ev event.Event) (*RoomshipCreatedEvent, error) {
	if ev.Topic() != TopicRoomshipCreated {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &RoomshipCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomshipCreatedEvent(
	roomshipID RoomshipID,
	generator event.IDGenerator,
) (*RoomshipCreatedEvent, error) {
	e := &RoomshipCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(roomshipID), TopicRoomshipCreated, payload, generator)
	return e, nil
}

func (e *RoomshipCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *RoomshipCreatedEvent) Unmarshal(_ []byte) error {
	return nil
}
