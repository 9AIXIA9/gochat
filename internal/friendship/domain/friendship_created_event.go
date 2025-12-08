package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendshipCreated event.Topic = "friendship.friendship.created"

var _ event.SpecificEvent = (*FriendshipCreatedEvent)(nil)

type FriendshipCreatedEvent struct {
	*event.StandardEvent
}

func ToFriendshipCreatedEvent(ev event.Event) (*FriendshipCreatedEvent, error) {
	if ev.Topic() != TopicFriendshipCreated {
		return nil, myErrors.ErrWrongEventTopic
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
	generator event.IDGenerator,
) (*FriendshipCreatedEvent, error) {
	e := &FriendshipCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(id), TopicFriendshipCreated, payload, generator)
	return e, nil
}

func (e *FriendshipCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *FriendshipCreatedEvent) Unmarshal([]byte) error {
	return nil
}
