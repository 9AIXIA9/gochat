package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendshipCreated event.Topic = "chat.friendship.created"

var _ event.SpecificEvent = (*FriendshipCreatedEvent)(nil)

type FriendshipCreatedEvent struct {
	userID1 kernel.UserID
	userID2 kernel.UserID
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
	id FriendshipID,
	userID1 kernel.UserID,
	userID2 kernel.UserID,
	generator event.IDGenerator,
) (*FriendshipCreatedEvent, error) {
	e := &FriendshipCreatedEvent{
		userID1: userID1,
		userID2: userID2,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(id), TopicFriendshipCreated, payload, generator)
	return e, nil
}

func (e *FriendshipCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		UserID1 kernel.UserID
		UserID2 kernel.UserID
	}
	return json.Marshal(Alias{
		UserID1: e.userID1,
		UserID2: e.userID2,
	})
}

func (e *FriendshipCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		UserID1 kernel.UserID
		UserID2 kernel.UserID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.userID1 = tmp.UserID1
	e.userID2 = tmp.UserID2
	return nil
}

func (e *FriendshipCreatedEvent) UserID1() kernel.UserID {
	return e.userID1
}

func (e *FriendshipCreatedEvent) UserID2() kernel.UserID {
	return e.userID2
}
