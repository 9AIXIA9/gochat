package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

const TopicReceivedRoomMessageNotificationRequested event.Topic = "notification.received_room_message_notification.requested"

var _ event.SpecificEvent = (*ReceivedRoomMessageNotificationRequestedEvent)(nil)

type ReceivedRoomMessageNotificationRequestedEvent struct {
	*event.StandardEvent
}

func ToReceivedRoomMessageNotificationRequestedEvent(ev event.Event) (*ReceivedRoomMessageNotificationRequestedEvent, error) {
	e := &ReceivedRoomMessageNotificationRequestedEvent{StandardEvent: event.NewStandardEventFrom(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func NewReceivedRoomMessageNotificationRequestedEvent(id event.ID, userID kernel.UserID) (*ReceivedRoomMessageNotificationRequestedEvent, error) {
	e := &ReceivedRoomMessageNotificationRequestedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(id, kernel.ID(userID), time.Now().UTC(), TopicReceivedRoomMessageNotificationRequested, payload)
	return e, nil
}

func (e *ReceivedRoomMessageNotificationRequestedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *ReceivedRoomMessageNotificationRequestedEvent) Unmarshal([]byte) error {
	return nil
}
