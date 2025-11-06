package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

//TODO 确定好所有事件的聚合ID

const TopicRoomJoined event.Topic = "social.room.joined"

var _ event.SpecificEvent = (*RoomJoinedEvent)(nil)

type RoomJoinedEvent struct {
	roomID RoomID
	userID kernel.UserID
	*event.StandardEvent
}

func ToRoomJoinedEvent(ev event.Event) (*RoomJoinedEvent, error) {
	e := &RoomJoinedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewRoomJoinedEvent(id event.ID, userID kernel.UserID, roomID RoomID) (*RoomJoinedEvent, error) {
	e := &RoomJoinedEvent{
		roomID: roomID,
		userID: userID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(roomID), time.Now().UTC(), TopicRoomJoined, payload)
	return e, nil
}

func (e *RoomJoinedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		RoomID RoomID
		UserID kernel.UserID
	}
	return json.Marshal(Alias{
		RoomID: e.roomID,
		UserID: e.userID,
	})
}

func (e *RoomJoinedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		RoomID RoomID
		UserID kernel.UserID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.roomID = tmp.RoomID
	e.userID = tmp.UserID
	return nil
}
