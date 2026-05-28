package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"time"
)

const TopicRoomCreated event.Topic = "room.created"

var _ event.SpecificEvent = (*RoomCreatedEvent)(nil)

type RoomCreatedEvent struct {
	createdAt time.Time
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
	type Alias struct {
		CreatedAt time.Time
	}
	return json.Marshal(Alias{
		CreatedAt: e.createdAt,
	})
}

func (e *RoomCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		CreatedAt time.Time
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.createdAt = tmp.CreatedAt
	return nil
}

func (e *RoomCreatedEvent) CreatedAt() time.Time {
	return e.createdAt
}
