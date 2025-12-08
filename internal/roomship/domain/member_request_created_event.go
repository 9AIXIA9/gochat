package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicMemberRequestCreated event.Topic = "roomship.member_request.created"

var _ event.SpecificEvent = (*MemberRequestCreatedEvent)(nil)

type MemberRequestCreatedEvent struct {
	*event.StandardEvent
}

func ToMemberRequestCreatedEvent(ev event.Event) (*MemberRequestCreatedEvent, error) {
	if ev.Topic() != TopicMemberRequestCreated {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &MemberRequestCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewMemberRequestCreatedEvent(
	requestID kernel.OperationID,
	generator event.IDGenerator,
) (*MemberRequestCreatedEvent, error) {
	e := &MemberRequestCreatedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(requestID), TopicMemberRequestCreated, payload, generator)
	return e, nil
}

func (e *MemberRequestCreatedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *MemberRequestCreatedEvent) Unmarshal(_ []byte) error {
	return nil
}
