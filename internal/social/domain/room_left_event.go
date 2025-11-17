package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicRoomLeft event.Topic = "social.room.left"

var _ event.SpecificEvent = (*RoomLeftEvent)(nil)

type RoomLeftEvent struct {
	userID kernel.UserID
	*event.StandardEvent
}

func ToRoomLeftEvent(ev event.Event) (*RoomLeftEvent, error) {
	e := &RoomLeftEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomLeftEvent(id event.ID, userID kernel.UserID, roomID kernel.RoomID) (*RoomLeftEvent, error) {
	e := &RoomLeftEvent{
		userID: userID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(roomID), time.Now().UTC(), TopicRoomLeft, payload)
	return e, nil
}

func (e *RoomLeftEvent) Marshal() ([]byte, error) {
	type Alias struct {
		UserID kernel.UserID
	}
	return json.Marshal(Alias{
		UserID: e.userID,
	})
}

func (e *RoomLeftEvent) Unmarshal(data []byte) error {
	type Alias struct {
		UserID kernel.UserID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.userID = tmp.UserID
	return nil
}

func (e *RoomLeftEvent) UserID() kernel.UserID {
	return e.userID
}
