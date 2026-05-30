package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendRequestAgreed event.Topic = "friend_request.agreed"

var _ event.SpecificEvent = (*FriendRequestAgreedEvent)(nil)

type FriendRequestAgreedEvent struct {
	*event.StandardEvent
}

func ToFriendRequestAgreedEvent(ev event.Event) (*FriendRequestAgreedEvent, error) {
	if ev.Topic() != TopicFriendRequestAgreed {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &FriendRequestAgreedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewFriendRequestAgreedEvent(
	requestID kernel.OperationID,
	generator event.IDGenerator,
) (*FriendRequestAgreedEvent, error) {
	e := &FriendRequestAgreedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(requestID), TopicFriendRequestAgreed, payload, generator)
	return e, nil
}

func (e *FriendRequestAgreedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *FriendRequestAgreedEvent) Unmarshal([]byte) error {
	return nil
}
