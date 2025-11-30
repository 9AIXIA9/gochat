package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendshipCreated event.Topic = "friendship.friendship.created"

var _ event.SpecificEvent = (*FriendshipCreatedEvent)(nil)

type FriendshipCreatedEvent struct {
	friendID kernel.UserID
	*event.StandardEvent
}

func ToFriendshipCreatedEvent(ev event.Event) (*FriendshipCreatedEvent, error) {
	if ev.Topic() != TopicFriendshipCreated {
		return nil, myErrors.ErrWrongEventType
	}
	e := &FriendshipCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewFriendshipCreatedEvent(
	userID kernel.UserID,
	friendID kernel.UserID,
	generator event.IDGenerator,
) (*FriendshipCreatedEvent, error) {
	e := &FriendshipCreatedEvent{
		friendID: friendID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicFriendshipCreated, payload, generator)
	return e, nil
}

func (e *FriendshipCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		FriendID kernel.UserID
	}
	return json.Marshal(&Alias{
		FriendID: e.friendID,
	})
}

func (e *FriendshipCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		FriendID kernel.UserID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.friendID = tmp.FriendID
	return nil
}

func (e *FriendshipCreatedEvent) FriendID() kernel.UserID {
	return e.friendID
}
