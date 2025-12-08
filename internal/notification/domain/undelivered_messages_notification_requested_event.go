package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicUndeliveredMessagesNotificationRequested event.Topic = "notification.undelivered_messages_notification.requested"

var _ event.SpecificEvent = (*UndeliveredMessagesNotificationRequestedEvent)(nil)

type UndeliveredMessagesNotificationRequestedEvent struct {
	*event.StandardEvent
}

func ToUndeliveredMessagesNotificationRequestedEvent(ev event.Event) (*UndeliveredMessagesNotificationRequestedEvent, error) {
	if ev.Topic() != TopicUndeliveredMessagesNotificationRequested {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &UndeliveredMessagesNotificationRequestedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUndeliveredMessagesNotificationRequestedEvent(
	userID kernel.UserID,
	generator event.IDGenerator,
) (*UndeliveredMessagesNotificationRequestedEvent, error) {
	e := &UndeliveredMessagesNotificationRequestedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(userID), TopicUndeliveredMessagesNotificationRequested, payload, generator)
	return e, nil
}

func (e *UndeliveredMessagesNotificationRequestedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UndeliveredMessagesNotificationRequestedEvent) Unmarshal([]byte) error {
	return nil
}
