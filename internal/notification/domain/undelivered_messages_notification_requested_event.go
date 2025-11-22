package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicUndeliveredMessagesNotificationRequested event.Topic = "notification.undelivered_messages_notification.requested"

var _ event.SpecificEvent = (*UndeliveredMessagesNotificationRequestedEvent)(nil)

type UndeliveredMessagesNotificationRequestedEvent struct {
	*event.StandardEvent
}

func ToUndeliveredMessagesNotificationRequestedEvent(ev event.Event) (*UndeliveredMessagesNotificationRequestedEvent, error) {
	e := &UndeliveredMessagesNotificationRequestedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUndeliveredMessagesNotificationRequestedEvent(
	id event.ID,
	userID kernel.UserID,
) (*UndeliveredMessagesNotificationRequestedEvent, error) {
	e := &UndeliveredMessagesNotificationRequestedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicUndeliveredMessagesNotificationRequested, payload)
	return e, nil
}

func (e *UndeliveredMessagesNotificationRequestedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UndeliveredMessagesNotificationRequestedEvent) Unmarshal([]byte) error {
	return nil
}
