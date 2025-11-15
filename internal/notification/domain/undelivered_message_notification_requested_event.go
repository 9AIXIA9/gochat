package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicUndeliveredMessageNotificationRequested event.Topic = "notification.undelivered_message_notification.requested"

var _ event.SpecificEvent = (*UndeliveredMessageNotificationRequestedEvent)(nil)

type UndeliveredMessageNotificationRequestedEvent struct {
	*event.StandardEvent
}

func ToUndeliveredMessageNotificationRequestedEvent(ev event.Event) (*UndeliveredMessageNotificationRequestedEvent, error) {
	e := &UndeliveredMessageNotificationRequestedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewUndeliveredMessageNotificationRequestedEvent(id event.ID, userID kernel.UserID) (*UndeliveredMessageNotificationRequestedEvent, error) {
	e := &UndeliveredMessageNotificationRequestedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicUndeliveredMessageNotificationRequested, payload)
	return e, nil
}

func (e *UndeliveredMessageNotificationRequestedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *UndeliveredMessageNotificationRequestedEvent) Unmarshal([]byte) error {
	return nil
}
