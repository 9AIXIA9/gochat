package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendRequestRefused event.Topic = "friendship.friend_request.refused"

var _ event.SpecificEvent = (*FriendRequestRefusedEvent)(nil)

type FriendRequestRefusedEvent struct {
	*event.StandardEvent
}

func ToFriendRequestRefusedEvent(ev event.Event) (*FriendRequestRefusedEvent, error) {
	if ev.Topic() != TopicFriendRequestRefused {
		return nil, myErrors.ErrWrongEventType
	}
	e := &FriendRequestRefusedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewFriendRequestRefusedEvent(
	requestID kernel.OperationID,
	generator event.IDGenerator,
) (*FriendRequestRefusedEvent, error) {
	e := &FriendRequestRefusedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(requestID), TopicFriendRequestRefused, payload, generator)
	return e, nil
}

func (e *FriendRequestRefusedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *FriendRequestRefusedEvent) Unmarshal([]byte) error {
	return nil
}
