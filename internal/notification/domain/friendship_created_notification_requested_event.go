package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendshipCreatedNotificationRequested event.Topic = "notification.friendship_created_notification.requested"

var _ event.SpecificEvent = (*FriendshipCreatedNotificationRequestedEvent)(nil)

type FriendshipCreatedNotificationRequestedEvent struct {
	friendID kernel.UserID
	*event.StandardEvent
}

func ToFriendshipCreatedNotificationRequestedEvent(ev event.Event) (*FriendshipCreatedNotificationRequestedEvent, error) {
	if ev.Topic() != TopicFriendshipCreatedNotificationRequested {
		return nil, myErrors.ErrWrongEventType
	}
	e := &FriendshipCreatedNotificationRequestedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewFriendshipCreatedNotificationRequestedEvent(
	userID kernel.UserID,
	friendID kernel.UserID,
	generator event.IDGenerator,
) (*FriendshipCreatedNotificationRequestedEvent, error) {
	e := &FriendshipCreatedNotificationRequestedEvent{
		friendID: friendID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(
		kernel.ID(userID),
		TopicFriendshipCreatedNotificationRequested,
		payload,
		generator,
	)
	return e, nil
}

func (e *FriendshipCreatedNotificationRequestedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		FriendID kernel.UserID
	}

	return json.Marshal(Alias{
		FriendID: e.friendID,
	})
}

func (e *FriendshipCreatedNotificationRequestedEvent) Unmarshal(data []byte) error {
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

func (e *FriendshipCreatedNotificationRequestedEvent) FriendID() kernel.UserID {
	return e.friendID
}
