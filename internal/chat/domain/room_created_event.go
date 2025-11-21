package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicRoomCreated event.Topic = "chat.room.created"

var _ event.SpecificEvent = (*RoomCreatedEvent)(nil)

type RoomCreatedEvent struct {
	ownerID kernel.UserID
	number  kernel.RoomNumber
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

func NewRoomCreatedEvent(id event.ID, roomID kernel.RoomID, number kernel.RoomNumber, ownerID kernel.UserID) (*RoomCreatedEvent, error) {
	e := &RoomCreatedEvent{
		ownerID: ownerID,
		number:  number,
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
		OwnerID kernel.UserID
		Number  kernel.RoomNumber
	}
	return json.Marshal(&Alias{
		OwnerID: e.ownerID,
		Number:  e.number,
	})
}

func (e *RoomCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		OwnerID kernel.UserID
		Number  kernel.RoomNumber
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.number = tmp.Number
	e.ownerID = tmp.OwnerID
	return nil
}

func (e *RoomCreatedEvent) Number() kernel.RoomNumber {
	return e.number
}

func (e *RoomCreatedEvent) OwnerID() kernel.UserID {
	return e.ownerID
}
