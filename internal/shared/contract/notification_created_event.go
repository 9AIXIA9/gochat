package contract

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicNotificationCreated event.Topic = "notification.created"

var _ event.SpecificEvent = (*NotificationCreatedEvent)(nil)

type NotificationCreatedEvent struct {
	*event.StandardEvent
	recipientID kernel.UserID
	rawPayload  json.RawMessage
}

func NewNotificationCreatedEvent(
	messageID kernel.MessageID,
	recipientID kernel.UserID,
	rawPayload json.RawMessage,
	generator event.IDGenerator,
) (*NotificationCreatedEvent, error) {
	e := &NotificationCreatedEvent{
		recipientID: recipientID,
		rawPayload:  rawPayload,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(messageID), TopicNotificationCreated, payload, generator)
	return e, nil
}

func ToNotificationCreatedEvent(ev event.Event) (*NotificationCreatedEvent, error) {
	if ev.Topic() != TopicNotificationCreated {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &NotificationCreatedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *NotificationCreatedEvent) RecipientID() kernel.UserID {
	return e.recipientID
}

func (e *NotificationCreatedEvent) RawPayload() []byte {
	return e.rawPayload
}

func (e *NotificationCreatedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		RecipientID kernel.UserID
		RawPayload  json.RawMessage
	}
	return json.Marshal(&Alias{
		RecipientID: e.recipientID,
		RawPayload:  e.rawPayload,
	})
}

func (e *NotificationCreatedEvent) Unmarshal([]byte) error {
	type Alias struct {
		RecipientID kernel.UserID
		RawPayload  json.RawMessage
	}
	var tmp Alias
	if err := json.Unmarshal(e.Payload(), &tmp); err != nil {
		return err
	}
	e.recipientID = tmp.RecipientID
	e.rawPayload = tmp.RawPayload
	return nil
}
