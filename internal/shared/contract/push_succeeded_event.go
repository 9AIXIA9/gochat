package contract

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicPushSucceeded event.Topic = "push.succeeded"

var _ event.SpecificEvent = (*PushSucceededEvent)(nil)

type PushSucceededEvent struct {
	*event.StandardEvent
}

func NewPushSucceededEvent(
	messageID kernel.MessageID,
	generator event.IDGenerator,
) (*PushSucceededEvent, error) {
	e := &PushSucceededEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(messageID), TopicPushSucceeded, payload, generator)
	return e, nil
}

func ToPushSucceededEvent(ev event.Event) (*PushSucceededEvent, error) {
	if ev.Topic() != TopicPushSucceeded {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &PushSucceededEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *PushSucceededEvent) Marshal() ([]byte, error) {
	return nil, nil
}

func (e *PushSucceededEvent) Unmarshal([]byte) error {
	return nil
}
