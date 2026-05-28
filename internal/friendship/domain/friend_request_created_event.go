package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendRequestCreated event.Topic = "friend_request.created"

var _ event.SpecificEvent = (*FriendRequestCreatedEvent)(nil)

type FriendRequestCreatedEvent struct {
	*event.StandardEvent
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
