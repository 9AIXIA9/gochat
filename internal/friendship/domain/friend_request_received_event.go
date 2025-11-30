package domain

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicFriendRequestReceived event.Topic = "friendship.friend_request.received"

var _ event.SpecificEvent = (*FriendRequestReceivedEvent)(nil)

type FriendRequestReceivedEvent struct {
	requestID kernel.OperationID
	*event.StandardEvent
}

func ToFriendRequestReceivedEvent(ev event.Event) (*FriendRequestReceivedEvent, error) {
	if ev.Topic() != TopicFriendRequestReceived {
		return nil, myErrors.ErrWrongEventType
	}
	e := &FriendRequestReceivedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewFriendRequestReceivedEvent(
	userID kernel.UserID,
	requestID kernel.OperationID,
	generator event.IDGenerator,
) (*FriendRequestReceivedEvent, error) {
	e := &FriendRequestReceivedEvent{
		requestID: requestID,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicFriendRequestReceived, payload, generator)
	return e, nil
}

func (e *FriendRequestReceivedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		RequestID kernel.OperationID
	}
	return json.Marshal(&Alias{
		RequestID: e.requestID,
	})
}

func (e *FriendRequestReceivedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		RequestID kernel.OperationID
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.requestID = tmp.RequestID
	return nil
}

func (e *FriendRequestReceivedEvent) RequestID() kernel.OperationID {
	return e.requestID
}
