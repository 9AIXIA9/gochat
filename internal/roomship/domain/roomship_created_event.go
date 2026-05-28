package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicRoomshipCreated event.Topic = "roomship.created"

var _ event.SpecificEvent = (*RoomshipCreatedEvent)(nil)

type RoomshipCreatedEvent struct {
	userID kernel.UserID
	roomID kernel.RoomID
	*event.StandardEvent
}

func NewRoomshipCreatedEvent(
	userID kernel.UserID,
	roomID kernel.RoomID,
	roomshipID RoomshipID,
	generator event.IDGenerator,
) (*RoomshipCreatedEvent, error) {
	e := &RoomshipCreatedEvent{
		userID: userID,
		roomID: roomID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(roomshipID), TopicRoomshipCreated, payload, generator)
	return e, nil
}

func (e *RoomshipCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		UserID kernel.UserID
		RoomID kernel.RoomID
	}
	return json.Marshal(&Alias{
		UserID: e.userID,
		RoomID: e.roomID,
	})
}

func (e *RoomshipCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		UserID kernel.UserID
		RoomID kernel.RoomID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.userID = tmp.UserID
	e.roomID = tmp.RoomID
	return nil
}
