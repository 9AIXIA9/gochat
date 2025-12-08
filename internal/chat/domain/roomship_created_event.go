package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicRoomshipCreated event.Topic = "chat.roomship.created"

var _ event.SpecificEvent = (*RoomshipCreatedEvent)(nil)

type RoomshipCreatedEvent struct {
	userID kernel.UserID
	roomID kernel.RoomID
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
	id RoomshipID,
	userID kernel.UserID,
	roomID kernel.RoomID,
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

	e.StandardEvent = event.NewStandardEvent(kernel.ID(id), TopicRoomshipCreated, payload, generator)
	return e, nil
}

func (e *RoomshipCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		UserID kernel.UserID
		RoomID kernel.RoomID
	}
	return json.Marshal(Alias{
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

func (e *RoomshipCreatedEvent) UserID() kernel.UserID {
	return e.userID
}

func (e *RoomshipCreatedEvent) RoomID() kernel.RoomID {
	return e.roomID
}
