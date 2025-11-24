package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicRoomJoined event.Topic = "social.room.joined"

var _ event.SpecificEvent = (*RoomJoinedEvent)(nil)

type RoomJoinedEvent struct {
	userID kernel.UserID
	*event.StandardEvent
}

func ToRoomJoinedEvent(ev event.Event) (*RoomJoinedEvent, error) {
	e := &RoomJoinedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomJoinedEvent(
	userID kernel.UserID,
	roomID kernel.RoomID,
	generator event.IDGenerator,
) (*RoomJoinedEvent, error) {
	e := &RoomJoinedEvent{
		userID: userID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(roomID), TopicRoomJoined, payload, generator)
	return e, nil
}

func (e *RoomJoinedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		UserID kernel.UserID
	}
	return json.Marshal(Alias{
		UserID: e.userID,
	})
}

func (e *RoomJoinedEvent) Unmarshal(data []byte) error {
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

func (e *RoomJoinedEvent) UserID() kernel.UserID {
	return e.userID
}
