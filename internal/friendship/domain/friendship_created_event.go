package domain

import (
	"encoding/json"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendshipCreated event.Topic = "friendship.created"

var _ event.SpecificEvent = (*FriendshipCreatedEvent)(nil)

type FriendshipCreatedEvent struct {
	userID1, userID2 kernel.UserID
	*event.StandardEvent
}

func NewFriendshipCreatedEvent(
	id FriendshipID,
	userID1, userID2 kernel.UserID,
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
		UserID1, UserID2 kernel.UserID
	}
	return json.Marshal(&Alias{
		UserID1: e.userID1,
		UserID2: e.userID2,
	})
}

func (e *FriendshipCreatedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		UserID1, UserID2 kernel.UserID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.userID1 = tmp.UserID1
	e.userID2 = tmp.UserID2
	return nil
}
