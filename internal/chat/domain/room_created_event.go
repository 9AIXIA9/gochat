package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicRoomCreated event.Topic = "chat.room.created"

var _ event.SpecificEvent = (*RoomCreatedEvent)(nil)

type RoomCreatedEvent struct {
	ownerID kernel.UserID
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
	ownerID kernel.UserID,
	generator event.IDGenerator,
) (*RoomCreatedEvent, error) {
	e := &RoomCreatedEvent{
		ownerID: ownerID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(roomID), TopicRoomCreated, payload, generator)
	return e, nil
}

func (e *RoomCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		OwnerID kernel.UserID
	}
	return json.Marshal(&Alias{
		OwnerID: e.ownerID,
	})
}

func (e *RoomCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		OwnerID kernel.UserID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.ownerID = tmp.OwnerID
	return nil
}

func (e *RoomCreatedEvent) OwnerID() kernel.UserID {
	return e.ownerID
}
