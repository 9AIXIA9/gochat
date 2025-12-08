package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendRequestCreated event.Topic = "friendship.friend_request.created"

var _ event.SpecificEvent = (*FriendRequestCreatedEvent)(nil)

type FriendRequestCreatedEvent struct {
	*event.StandardEvent
}

func ToFriendRequestCreatedEvent(ev event.Event) (*FriendRequestCreatedEvent, error) {
	if ev.Topic() != TopicFriendRequestCreated {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &FriendRequestCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewFriendRequestCreatedEvent(
	requestID kernel.OperationID,
	generator event.IDGenerator,
) (*FriendRequestCreatedEvent, error) {
	e := &FriendRequestCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(requestID), TopicFriendRequestCreated, payload, generator)
	return e, nil
}

func (e *FriendRequestCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *FriendRequestCreatedEvent) Unmarshal([]byte) error {
	return nil
}
