package contract

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

const TopicPushRequested event.Topic = "push.requested"

var _ event.SpecificEvent = (*PushRequestedEvent)(nil)

type PushRequestedEvent struct {
	*event.StandardEvent
	recipientID kernel.UserID
	rawPayload  json.RawMessage
}

func NewPushRequestedEvent(
	messageID kernel.MessageID,
	recipientID kernel.UserID,
	rawPayload json.RawMessage,
	generator event.IDGenerator,
) (*PushRequestedEvent, error) {
	e := &PushRequestedEvent{
		recipientID: recipientID,
		rawPayload:  rawPayload,
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = event.NewStandardEvent(kernel.ID(messageID), TopicPushRequested, payload, generator)
	return e, nil
}

func ToPushRequestedEvent(ev event.Event) (*PushRequestedEvent, error) {
	if ev.Topic() != TopicPushRequested {
		return nil, myErrors.ErrWrongEventTopic
	}
	e := &PushRequestedEvent{StandardEvent: event.LoadStandardEventFromEvent(ev)}
	if len(ev.Payload()) > 0 {
		if err := e.Unmarshal(ev.Payload()); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *PushRequestedEvent) RecipientID() kernel.UserID {
	return e.recipientID
}

func (e *PushRequestedEvent) RawPayload() []byte {
	return e.rawPayload
}

func (e *PushRequestedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		RecipientID kernel.UserID
		RawPayload  json.RawMessage
	}
	return json.Marshal(&Alias{
		RecipientID: e.recipientID,
		RawPayload:  e.rawPayload,
	})
}

func (e *PushRequestedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		RecipientID kernel.UserID
		RawPayload  json.RawMessage
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.recipientID = tmp.RecipientID
	e.rawPayload = tmp.RawPayload
	return nil
}
