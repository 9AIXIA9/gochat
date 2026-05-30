package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicMemberRequestAgreed event.Topic = "member_request.agreed"

var _ event.SpecificEvent = (*MemberRequestAgreedEvent)(nil)

type MemberRequestAgreedEvent struct {
	*event.StandardEvent
}

func ToMemberRequestAgreedEvent(ev event.Event) (*MemberRequestAgreedEvent, error) {
	if ev.Topic() != TopicMemberRequestAgreed {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &MemberRequestAgreedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewMemberRequestAgreedEvent(
	requestID kernel.OperationID,
	generator event.IDGenerator,
) (*MemberRequestAgreedEvent, error) {
	e := &MemberRequestAgreedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(requestID), TopicMemberRequestAgreed, payload, generator)
	return e, nil
}

func (e *MemberRequestAgreedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *MemberRequestAgreedEvent) Unmarshal(_ []byte) error {
	return nil
}
