package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicReceivedPrivateMessageNotificationRequested event.Topic = "notification.received_private_message_notification.requested"

var _ event.SpecificEvent = (*ReceivedPrivateMessageNotificationRequestedEvent)(nil)

type ReceivedPrivateMessageNotificationRequestedEvent struct {
	*event.StandardEvent
}

func ToReceivedPrivateMessageNotificationRequestedEvent(ev event.Event) (*ReceivedPrivateMessageNotificationRequestedEvent, error) {
	e := &ReceivedPrivateMessageNotificationRequestedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewReceivedPrivateMessageNotificationRequestedEvent(id event.ID, userID kernel.UserID) (*ReceivedPrivateMessageNotificationRequestedEvent, error) {
	e := &ReceivedPrivateMessageNotificationRequestedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicReceivedPrivateMessageNotificationRequested, payload)
	return e, nil
}

func (e *ReceivedPrivateMessageNotificationRequestedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *ReceivedPrivateMessageNotificationRequestedEvent) Unmarshal([]byte) error {
	return nil
}
